package enumerator

import (
	"errors"
	"os"
	"reflect"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/khulnasoft-lab/driftctl/pkg/iac/config"
	awstest "github.com/khulnasoft-lab/driftctl/test/aws"
	"github.com/stretchr/testify/mock"
)

func TestS3Enumerator_NewS3Enumerator(t *testing.T) {
	tests := []struct {
		name   string
		config config.SupplierConfig
		setEnv map[string]string
		want   string
	}{
		{
			name: "test with no proxy env var",
			config: config.SupplierConfig{
				Key:     "tfstate",
				Backend: "s3",
				Path:    "terraform.tfstate",
			},
			setEnv: map[string]string{
				"AWS_DEFAULT_REGION": "us-east-1",
			},
			want: "us-east-1",
		},
		{
			name: "test with proxy env var",
			config: config.SupplierConfig{
				Key:     "tfstate",
				Backend: "s3",
				Path:    "terraform.tfstate",
			},
			setEnv: map[string]string{
				"AWS_DEFAULT_REGION":     "us-east-1",
				"DCTL_S3_DEFAULT_REGION": "eu-west-3",
			},
			want: "eu-west-3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.setEnv {
				os.Setenv(key, value)
			}
			got := NewS3Enumerator(tt.config).client.(*s3.S3).Config.Region
			if awssdk.StringValue(got) != tt.want {
				t.Errorf("NewS3Enumerator().client.Config.Region got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS3Enumerator_Enumerate(t *testing.T) {
	tests := []struct {
		name   string
		config config.SupplierConfig
		mocks  func(client *awstest.MockFakeS3)
		want   []string
		err    string
	}{
		{
			name: "no test results are returned",
			config: config.SupplierConfig{
				Path: "bucket-name/a/nested/prefix",
			},
			mocks: func(client *awstest.MockFakeS3) {
				input := &s3.ListObjectsV2Input{
					Bucket: awssdk.String("bucket-name"),
					Prefix: awssdk.String("a/nested/prefix"),
				}
				client.On(
					"ListObjectsV2Pages",
					input,
					mock.MatchedBy(func(callback func(res *s3.ListObjectsV2Output, lastPage bool) bool) bool {
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/state1"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/state2"),
									Size: awssdk.Int64(2),
								},
								{
									Key:  awssdk.String("a/nested/prefix/state3"),
									Size: awssdk.Int64(1),
								},
							},
						}, false)
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/state4"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/folder1/state5"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/folder2/subfolder1/state6"),
									Size: awssdk.Int64(5),
								},
							},
						}, true)
						return true
					}),
				).Return(nil)
			},
			want: []string{},
			err:  "no Terraform state was found in bucket-name/a/nested/prefix, exiting",
		},
		{
			name: "one test result is returned",
			config: config.SupplierConfig{
				Path: "bucket-name/a/nested/prefix/state2",
			},
			mocks: func(client *awstest.MockFakeS3) {
				input := &s3.ListObjectsV2Input{
					Bucket: awssdk.String("bucket-name"),
					Prefix: awssdk.String("a/nested/prefix/state2"),
				}
				client.On(
					"ListObjectsV2Pages",
					input,
					mock.MatchedBy(func(callback func(res *s3.ListObjectsV2Output, lastPage bool) bool) bool {
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/state1"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/state2"),
									Size: awssdk.Int64(2),
								},
								{
									Key:  awssdk.String("a/nested/prefix/state3"),
									Size: awssdk.Int64(1),
								},
							},
						}, false)
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/state4"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/folder1/state5"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/folder2/subfolder1/state6"),
									Size: awssdk.Int64(5),
								},
							},
						}, true)
						return true
					}),
				).Return(nil)
			},
			want: []string{"bucket-name/a/nested/prefix/state2"},
		},
		{
			name: "test results with simple doublestar glob",
			config: config.SupplierConfig{
				Path: "bucket-name/**/*.tfstate",
			},
			mocks: func(client *awstest.MockFakeS3) {
				input := &s3.ListObjectsV2Input{
					Bucket: awssdk.String("bucket-name"),
					Prefix: awssdk.String(""),
				}
				client.On(
					"ListObjectsV2Pages",
					input,
					mock.MatchedBy(func(callback func(res *s3.ListObjectsV2Output, lastPage bool) bool) bool {
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/1/state1.tfstate"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/folder1/2/state2.tfstate"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/state3.tfstate"),
									Size: awssdk.Int64(5),
								},
							},
						}, false)
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/4/4/state4.tfstate"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/state5.tfstate"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/state6.tfstate.backup"),
									Size: awssdk.Int64(5),
								},
							},
						}, true)
						return true
					}),
				).Return(nil)
			},
			want: []string{
				"bucket-name/a/nested/prefix/1/state1.tfstate",
				"bucket-name/a/nested/folder1/2/state2.tfstate",
				"bucket-name/a/nested/prefix/state3.tfstate",
				"bucket-name/a/nested/prefix/4/4/state4.tfstate",
				"bucket-name/a/nested/state5.tfstate",
			},
			err: "",
		},
		{
			name: "test results with glob and prefix after glob",
			config: config.SupplierConfig{
				Path: "bucket-name/a/**/b/*.tfstate",
			},
			mocks: func(client *awstest.MockFakeS3) {
				input := &s3.ListObjectsV2Input{
					Bucket: awssdk.String("bucket-name"),
					Prefix: awssdk.String("a"),
				}
				client.On(
					"ListObjectsV2Pages",
					input,
					mock.MatchedBy(func(callback func(res *s3.ListObjectsV2Output, lastPage bool) bool) bool {
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/prefix/b/state1.tfstate"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/b/state2.tfstate"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/prefix/state3.tfstate"),
									Size: awssdk.Int64(5),
								}, {
									Key:  awssdk.String("a/prefix/state4.tfstate.backup"),
									Size: awssdk.Int64(5),
								},
							},
						}, true)
						return true
					}),
				).Return(nil)
			},
			want: []string{
				"bucket-name/a/prefix/b/state1.tfstate",
				"bucket-name/a/b/state2.tfstate",
			},
			err: "",
		},
		{
			name: "test results with glob",
			config: config.SupplierConfig{
				Path: "bucket-name/a/nested/prefix/**/*.tfstate",
			},
			mocks: func(client *awstest.MockFakeS3) {
				input := &s3.ListObjectsV2Input{
					Bucket: awssdk.String("bucket-name"),
					Prefix: awssdk.String("a/nested/prefix"),
				}
				client.On(
					"ListObjectsV2Pages",
					input,
					mock.MatchedBy(func(callback func(res *s3.ListObjectsV2Output, lastPage bool) bool) bool {
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/1/state1.tfstate"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/folder1/2/state2.tfstate"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/state3.tfstate"),
									Size: awssdk.Int64(5),
								},
							},
						}, false)
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/4/4/state4.tfstate"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/state5.state"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/state6.tfstate.backup"),
									Size: awssdk.Int64(5),
								},
							},
						}, true)
						return true
					}),
				).Return(nil)
			},
			want: []string{
				"bucket-name/a/nested/prefix/1/state1.tfstate",
				"bucket-name/a/nested/prefix/state3.tfstate",
				"bucket-name/a/nested/prefix/4/4/state4.tfstate",
			},
			err: "",
		},
		{
			name: "test results with simple glob",
			config: config.SupplierConfig{
				Path: "bucket-name/a/nested/prefix/*.tfstate",
			},
			mocks: func(client *awstest.MockFakeS3) {
				input := &s3.ListObjectsV2Input{
					Bucket: awssdk.String("bucket-name"),
					Prefix: awssdk.String("a/nested/prefix"),
				}
				client.On(
					"ListObjectsV2Pages",
					input,
					mock.MatchedBy(func(callback func(res *s3.ListObjectsV2Output, lastPage bool) bool) bool {
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/1/state1.tfstate"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/2/state2.tfstate"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/state3.tfstate"),
									Size: awssdk.Int64(5),
								},
							},
						}, false)
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/4/4/state4.tfstate"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/state5.state"),
									Size: awssdk.Int64(5),
								},
								{
									Key:  awssdk.String("a/nested/prefix/state6.tfstate.backup"),
									Size: awssdk.Int64(5),
								},
							},
						}, true)
						return true
					}),
				).Return(nil)
			},
			want: []string{"bucket-name/a/nested/prefix/state3.tfstate"},
			err:  "",
		},
		{
			name: "test when invalid config used",
			config: config.SupplierConfig{
				Path: "bucket-name",
			},
			mocks: func(client *awstest.MockFakeS3) {
				client.On("ListObjectsV2Pages", mock.Anything, mock.Anything).Return(errors.New("error when listing"))
			},
			want: nil,
			err:  "Unable to parse S3 path: bucket-name. Must be BUCKET_NAME/PREFIX",
		},
		{
			name:   "test when empty config used",
			config: config.SupplierConfig{},
			mocks: func(client *awstest.MockFakeS3) {
				client.On("ListObjectsV2Pages", mock.Anything, mock.Anything).Return(errors.New("error when listing"))
			},
			want: nil,
			err:  "Unable to parse S3 path: . Must be BUCKET_NAME/PREFIX",
		},
		{
			name: "test enumeration return error",
			config: config.SupplierConfig{
				Path: "bucket-name/a/nested/prefix",
			},
			mocks: func(client *awstest.MockFakeS3) {
				client.On("ListObjectsV2Pages", mock.Anything, mock.Anything).Return(errors.New("error when listing"))
			},
			want: nil,
			err:  "error when listing",
		},
		{
			name: "test no state found with simple path",
			config: config.SupplierConfig{
				Path: "bucket-name/a/nested/prefix",
			},
			mocks: func(client *awstest.MockFakeS3) {
				input := &s3.ListObjectsV2Input{
					Bucket: awssdk.String("bucket-name"),
					Prefix: awssdk.String("a/nested/prefix"),
				}
				client.On(
					"ListObjectsV2Pages",
					input,
					mock.MatchedBy(func(callback func(res *s3.ListObjectsV2Output, lastPage bool) bool) bool {
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/1/state1.tfstate"),
									Size: awssdk.Int64(5),
								},
							},
						}, true)
						return true
					}),
				).Return(nil)
			},
			want: []string{},
			err:  "no Terraform state was found in bucket-name/a/nested/prefix, exiting",
		},
		{
			name: "test no state found with simple glob path",
			config: config.SupplierConfig{
				Path: "bucket-name/a/nested/prefix/*",
			},
			mocks: func(client *awstest.MockFakeS3) {
				input := &s3.ListObjectsV2Input{
					Bucket: awssdk.String("bucket-name"),
					Prefix: awssdk.String("a/nested/prefix"),
				}
				client.On(
					"ListObjectsV2Pages",
					input,
					mock.MatchedBy(func(callback func(res *s3.ListObjectsV2Output, lastPage bool) bool) bool {
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/1/state1.tfstate"),
									Size: awssdk.Int64(5),
								},
							},
						}, true)
						return true
					}),
				).Return(nil)
			},
			want: []string{},
			err:  "no Terraform state was found in bucket-name/a/nested/prefix/*, exiting",
		},
		{
			name: "test no state found with double star glob path",
			config: config.SupplierConfig{
				Path: "bucket-name/a/nested/prefix/**/*.tfstate",
			},
			mocks: func(client *awstest.MockFakeS3) {
				input := &s3.ListObjectsV2Input{
					Bucket: awssdk.String("bucket-name"),
					Prefix: awssdk.String("a/nested/prefix"),
				}
				client.On(
					"ListObjectsV2Pages",
					input,
					mock.MatchedBy(func(callback func(res *s3.ListObjectsV2Output, lastPage bool) bool) bool {
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/1/dummy.json"),
									Size: awssdk.Int64(5),
								},
							},
						}, true)
						return true
					}),
				).Return(nil)
			},
			want: []string{},
			err:  "no Terraform state was found in bucket-name/a/nested/prefix/**/*.tfstate, exiting",
		},
		{
			name: "test folder terraform.tfstate is not recognized as a file",
			config: config.SupplierConfig{
				Path: "bucket-name/a/nested/**/*.tfstate",
			},
			mocks: func(client *awstest.MockFakeS3) {
				input := &s3.ListObjectsV2Input{
					Bucket: awssdk.String("bucket-name"),
					Prefix: awssdk.String("a/nested"),
				}
				client.On(
					"ListObjectsV2Pages",
					input,
					mock.MatchedBy(func(callback func(res *s3.ListObjectsV2Output, lastPage bool) bool) bool {
						callback(&s3.ListObjectsV2Output{
							Contents: []*s3.Object{
								{
									Key:  awssdk.String("a/nested/prefix/terraform.tfstate/terraform.tfstate"),
									Size: awssdk.Int64(5),
								},
							},
						}, true)
						return true
					}),
				).Return(nil)
			},
			want: []string{"bucket-name/a/nested/prefix/terraform.tfstate/terraform.tfstate"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeS3 := awstest.MockFakeS3{}
			tt.mocks(&fakeS3)
			s := &S3Enumerator{
				config: tt.config,
				client: &fakeS3,
			}
			got, err := s.Enumerate()
			if err != nil && err.Error() != tt.err {
				t.Fatalf("Expected error '%s', got '%s'", tt.err, err.Error())
				return
			}
			if tt.err != "" && err == nil {
				t.Fatalf("Expected error '%s' but got nil", tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Enumerate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
