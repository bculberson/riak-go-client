package riak

import (
	"fmt"
	rpbRiakKV "github.com/basho-labs/riak-go-client/rpb/riak_kv"
	proto "github.com/golang/protobuf/proto"
	"reflect"
	"time"
)

// FetchValueCommand

type FetchValueCommandOptions struct {
	Location
	R                   uint32
	Pr                  uint32
	BasicQuorum         bool
	NotFoundOk          bool
	IfNotModified       []byte // TODO pb field is IfModified
	HeadOnly            bool
	ReturnDeletedVClock bool
	Timeout             time.Duration
	SloppyQuorum        bool
	NVal                uint32
	// TODO ConflictResolver
}

func (options *FetchValueCommandOptions) GetTimeoutMilliseconds() *uint32 {
	if options.Timeout > 0 {
		timeoutMilliseconds := uint32(options.Timeout / time.Millisecond)
		return &timeoutMilliseconds
	} else {
		return nil
	}
}

type FetchValueCommand struct {
	CommandImpl
	options  *FetchValueCommandOptions
	Response *FetchValueResponse
}

func NewFetchValueCommand(options *FetchValueCommandOptions) (cmd *FetchValueCommand, err error) {
	if options == nil {
		// TODO default options?
		err = ErrNilOptions
		return
	}

	if err = validateObjectLocator(options); err != nil {
		return
	}

	cmd = &FetchValueCommand{
		options: options,
	}

	return
}

func (cmd *FetchValueCommand) Name() string {
	return "FetchValue"
}

func (cmd *FetchValueCommand) constructPbRequest() (proto.Message, error) {
	rpb := &rpbRiakKV.RpbGetReq{
		Type:          []byte(cmd.options.BucketType),
		Bucket:        []byte(cmd.options.Bucket),
		Key:           []byte(cmd.options.Key),
		BasicQuorum:   &cmd.options.BasicQuorum,
		NotfoundOk:    &cmd.options.NotFoundOk,
		IfModified:    cmd.options.IfNotModified,
		Head:          &cmd.options.HeadOnly,
		Deletedvclock: &cmd.options.ReturnDeletedVClock,
		Timeout:       cmd.options.GetTimeoutMilliseconds(),
		SloppyQuorum:  &cmd.options.SloppyQuorum,
	}
	if cmd.options.R > 0 {
		rpb.R = &cmd.options.R
	}
	if cmd.options.Pr > 0 {
		rpb.Pr = &cmd.options.Pr
	}
	if cmd.options.NVal > 0 {
		rpb.NVal = &cmd.options.NVal
	}
	return rpb, nil
}

func (cmd *FetchValueCommand) onSuccess(msg proto.Message) error {
	if msg == nil {
		cmd.Response = &FetchValueResponse{
			IsNotFound:  true,
			IsUnchanged: false,
		}
	} else {
		if rpbGetResp, ok := msg.(*rpbRiakKV.RpbGetResp); ok {
			vclock := rpbGetResp.GetVclock()
			response := &FetchValueResponse{
				VClock:      vclock,
				IsUnchanged: rpbGetResp.GetUnchanged(),
				IsNotFound:  false,
			}

			if pbContent := rpbGetResp.GetContent(); pbContent == nil || len(pbContent) == 0 {
				object := &Object{
					IsTombstone: true,
					BucketType:  cmd.options.BucketType,
					Bucket:      cmd.options.Bucket,
					Key:         cmd.options.Key,
				}
				response.Values = []*Object{object}
			} else {
				response.Values = make([]*Object, len(pbContent))
				for i, content := range pbContent {
					if ro, err := NewObjectFromRpbContent(content); err != nil {
						return err
					} else {
						ro.VClock = vclock
						ro.BucketType = cmd.options.BucketType
						ro.Bucket = cmd.options.Bucket
						ro.Key = cmd.options.Key
						/*
							* TODO
							if (this.options.conflictResolver) {
								values = [this.options.conflictResolver(values)];
							}
						*/
						response.Values[i] = ro
					}
				}
			}

			cmd.Response = response
		} else {
			// TODO specific Riak error?
			return fmt.Errorf("[FetchValueCommand] could not convert %v to RpbGetResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *FetchValueCommand) getRequestCode() byte {
	return rpbCode_RpbGetReq
}

func (cmd *FetchValueCommand) getExpectedResponseCode() byte {
	return rpbCode_RpbGetResp
}

func (cmd *FetchValueCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakKV.RpbGetResp{}
}

// FetchValueResponse

type FetchValueResponse struct {
	IsNotFound  bool
	IsUnchanged bool
	VClock      []byte
	Values      []*Object
}

// FetchValueCommandBuilder

type FetchValueCommandBuilder struct {
	Options *FetchValueCommandOptions
}

func NewFetchValueCommandBuilder() *FetchValueCommandBuilder {
	builder := &FetchValueCommandBuilder{
		Options: &FetchValueCommandOptions{},
	}
	return builder
}

func (builder *FetchValueCommandBuilder) WithBucketType(bucketType string) *FetchValueCommandBuilder {
	builder.Options.BucketType = bucketType
	return builder
}

func (builder *FetchValueCommandBuilder) WithBucket(bucket string) *FetchValueCommandBuilder {
	builder.Options.Bucket = bucket
	return builder
}

func (builder *FetchValueCommandBuilder) WithKey(key string) *FetchValueCommandBuilder {
	builder.Options.Key = key
	return builder
}

func (builder *FetchValueCommandBuilder) Build() (Command, error) {
	if builder.Options == nil {
		return nil, ErrNilOptions
	}
	return NewFetchValueCommand(builder.Options)
}