package storage

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/iuliailies/photo-flux/internal/config"
	madmin "github.com/minio/madmin-go/v2"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/notification"
)

type Storage struct {
	conn           Conn
	userPolicyName string
}

// New initializes a storage object. The storage object should be shared accross multiple requests.
func New(config config.Storage) (*Storage, error) {

	endpoint := fmt.Sprintf("%s:%d", config.MinioAddress, config.MinioPort)

	admin, err := madmin.New(endpoint, config.AccessKey, config.SecretKey, false)
	if err != nil {
		return nil, fmt.Errorf("could not initialize admin connection to minio: %w", err)
	}

	opts := minio.Options{
		Secure: false,
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
	}

	minioClient, err := minio.New(endpoint, &opts)
	if err != nil {
		return nil, fmt.Errorf("could not initialize minio client with admin privilegies: %w", err)
	}

	return &Storage{
		conn: Conn{
			admin:       admin,
			adminClient: minioClient,
			adminLock:   &sync.Mutex{},
		},
		userPolicyName: config.UserPolicyName,
	}, nil
}

// NewMinioUser configures a user in the minio storage: add user + set policies + create corresponding bucket
func (s *Storage) NewMinioUser(ctx context.Context, id uuid.UUID, secret string) error {
	s.conn.adminLock.Lock()
	defer s.conn.adminLock.Unlock()

	userid := id.String()

	err := s.addUser(ctx, userid, secret)
	if err != nil {
		defer s.FailSafe(ctx, userid)
		return fmt.Errorf("could not create minio user: %w", err)
	}

	err = s.addBucket(ctx, userid, secret)
	if err != nil {
		defer s.FailSafe(ctx, userid)
		return fmt.Errorf("could not create minio bucket: %w", err)
	}

	return nil
}

// addUser adds a new user to the minio storage server
func (s *Storage) addUser(ctx context.Context, userid, secret string) error {
	err := s.conn.admin.AddUser(ctx, userid, secret)
	if err != nil {
		return fmt.Errorf("could not add user to minio: %w", err)
	}

	// Set this user's policy to a predefined policy.
	err = s.conn.admin.AttachPolicy(ctx, madmin.PolicyAssociationReq{
		User:     userid,
		Policies: []string{s.userPolicyName},
	})
	if err != nil {
		return fmt.Errorf("could not assign policy to user %s: %w", userid, err)
	}

	return nil
}

// addBucket adds a bucket for the specific user. Should always be called after addUser.
func (s *Storage) addBucket(ctx context.Context, userid, secret string) error {

	bucket := getBucketName(userid)

	err := s.conn.adminClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("could not create bucket for user %s: %w", userid, err)
	}

	// binds to queue, to send create and delete events to rabbitmq when they will happen
	queueArn := notification.NewArn("minio", "sqs", "", "PRIMARY", "amqp")
	queueConfig := notification.NewConfig(queueArn)
	queueConfig.AddEvents(notification.ObjectCreatedAll, notification.ObjectRemovedAll)

	notifconfig := notification.Configuration{}
	notifconfig.AddQueue(queueConfig)

	err = s.conn.adminClient.SetBucketNotification(ctx, bucket, notifconfig)
	if err != nil {
		return fmt.Errorf("could not enable notifications for bucket %s: %w", bucket, err)
	}

	err = s.conn.admin.SetBucketQuota(ctx, bucket, &madmin.BucketQuota{
		Type:  madmin.HardQuota,
		Quota: 5 * 1 << 30, // 5gb
	})
	if err != nil {
		return fmt.Errorf("could not set bucket quota for bucket %s: %w", bucket, err)
	}
	return nil
}

// FailSafe attepts to perform cleanup operation if the user creation process fails.
// Note that the admin connections should still be locked at this point.
func (s *Storage) FailSafe(ctx context.Context, userid string) {
	err := s.conn.adminClient.RemoveBucket(ctx, getBucketName(userid))

	// TODO: nice logs
	if err != nil {
		fmt.Println("removing bucket", err)
	}

	err = s.conn.admin.RemoveUser(ctx, userid)
	if err != nil {
		fmt.Println("removing user", err)
	}
}

func getBucketName(userid string) string {
	bucket := strings.Join([]string{"user", userid}, "-")
	return bucket
}
