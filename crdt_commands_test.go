package riak

import (
	"bytes"
	"fmt"
	rpbRiakDT "github.com/basho-labs/riak-go-client/rpb/riak_dt"
	"reflect"
	"testing"
	"time"
)

// UpdateCounter
// DtUpdateReq
// DtUpdateResp

func TestBuildDtUpdateReqCorrectlyViaUpdateCounterCommandBuilder(t *testing.T) {
	builder := NewUpdateCounterCommandBuilder().
		WithBucketType("counters").
		WithBucket("myBucket").
		WithKey("counter_1").
		WithIncrement(100).
		WithW(3).
		WithPw(1).
		WithDw(2).
		WithReturnBody(true).
		WithTimeout(time.Second * 20)
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}
	if req, ok := protobuf.(*rpbRiakDT.DtUpdateReq); ok {
		if expected, actual := "counters", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "myBucket", string(req.GetBucket()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "counter_1", string(req.GetKey()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(3), req.GetW(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(1), req.GetPw(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(2), req.GetDw(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetReturnBody(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		op := req.Op.CounterOp
		if expected, actual := int64(100), op.GetIncrement(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		validateTimeout(t, time.Second*20, req.GetTimeout())
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakDT.DtUpdateReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestUpdateCounterParsesDtUpdateRespCorrectly(t *testing.T) {
	counterValue := int64(1234)
	generatedKey := "generated_key"
	dtUpdateResp := &rpbRiakDT.DtUpdateResp{
		CounterValue: &counterValue,
		Key:          []byte(generatedKey),
	}

	builder := NewUpdateCounterCommandBuilder().
		WithBucketType("counters").
		WithBucket("myBucket").
		WithKey("counter_1").
		WithIncrement(100)
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}

	cmd.onSuccess(dtUpdateResp)

	if uc, ok := cmd.(*UpdateCounterCommand); ok {
		rsp := uc.Response
		if expected, actual := int64(1234), rsp.CounterValue; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "generated_key", rsp.GeneratedKey; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *UpdateCounterCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestValidationOfUpdateCounterViaBuilder(t *testing.T) {
	// validate that Bucket is required
	builder := NewUpdateCounterCommandBuilder()
	_, err := builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrBucketRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	// validate that Key is NOT required
	builder = NewUpdateCounterCommandBuilder()
	builder.WithBucket("bucket_name")
	_, err = builder.Build()
	if err != nil {
		t.Fatal("expected nil err")
	}
}

// FetchCounter
// DtFetchReq
// DtFetchResp

func TestBuildDtFetchReqCorrectlyViaFetchCounterCommandBuilder(t *testing.T) {
	builder := NewFetchCounterCommandBuilder().
		WithBucketType("counters").
		WithBucket("myBucket").
		WithKey("counter_1").
		WithR(3).
		WithPr(1).
		WithNotFoundOk(true).
		WithBasicQuorum(true).
		WithTimeout(time.Second * 20)
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}
	if req, ok := protobuf.(*rpbRiakDT.DtFetchReq); ok {
		if expected, actual := "counters", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "myBucket", string(req.GetBucket()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "counter_1", string(req.GetKey()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(3), req.GetR(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(1), req.GetPr(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetNotfoundOk(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetBasicQuorum(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		validateTimeout(t, time.Second*20, req.GetTimeout())
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakDT.DtFetchReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestFetchCounterParsesDtFetchRespCorrectly(t *testing.T) {
	counterValue := int64(1234)
	dtValue := &rpbRiakDT.DtValue{
		CounterValue: &counterValue,
	}
	dtFetchResp := &rpbRiakDT.DtFetchResp{
		Value: dtValue,
	}

	builder := NewFetchCounterCommandBuilder().
		WithBucketType("counters").
		WithBucket("myBucket").
		WithKey("counter_1")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}

	cmd.onSuccess(dtFetchResp)

	if uc, ok := cmd.(*FetchCounterCommand); ok {
		rsp := uc.Response
		if expected, actual := counterValue, rsp.CounterValue; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *FetchCounterCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestFetchCounterParsesDtFetchRespWithoutValueCorrectly(t *testing.T) {
	builder := NewFetchCounterCommandBuilder().
		WithBucketType("counters").
		WithBucket("myBucket").
		WithKey("counter_1")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}

	dtFetchResp := &rpbRiakDT.DtFetchResp{}
	cmd.onSuccess(dtFetchResp)

	if uc, ok := cmd.(*FetchCounterCommand); ok {
		if expected, actual := true, uc.Response.IsNotFound; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *FetchCounterCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestValidationOfFetchCounterViaBuilder(t *testing.T) {
	// validate that Bucket is required
	builder := NewFetchCounterCommandBuilder()
	_, err := builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrBucketRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	// validate that Key is required
	builder = NewFetchCounterCommandBuilder()
	builder.WithBucket("bucket_name")
	_, err = builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrKeyRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

// UpdateSet
// DtUpdateReq
// DtUpdateResp

func TestBuildDtUpdateReqCorrectlyViaUpdateSetCommandBuilder(t *testing.T) {
	builder := NewUpdateSetCommandBuilder().
		WithBucketType("sets").
		WithBucket("bucket").
		WithKey("key").
		WithContext(crdtContextBytes).
		WithAdditions([]byte("a1"), []byte("a2")).
		WithAdditions([]byte("a3"), []byte("a4")).
		WithRemovals([]byte("r1"), []byte("r2")).
		WithRemovals([]byte("r3"), []byte("r4")).
		WithW(1).
		WithDw(2).
		WithPw(3).
		WithReturnBody(true).
		WithTimeout(time.Second * 20)
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}
	if req, ok := protobuf.(*rpbRiakDT.DtUpdateReq); ok {
		if expected, actual := "sets", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "bucket", string(req.GetBucket()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "key", string(req.GetKey()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := 0, bytes.Compare(crdtContextBytes, req.GetContext()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(1), req.GetW(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(2), req.GetDw(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(3), req.GetPw(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		validateTimeout(t, time.Second*20, req.GetTimeout())

		op := req.Op.SetOp

		for i := 1; i <= 4; i++ {
			aitem := fmt.Sprintf("a%d", i)
			ritem := fmt.Sprintf("r%d", i)
			if expected, actual := true, sliceIncludes(op.Adds, []byte(aitem)); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := true, sliceIncludes(op.Removes, []byte(ritem)); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakDT.DtUpdateReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestUpdateSetParsesDtUpdateRespCorrectly(t *testing.T) {
	setValue := [][]byte{
		[]byte("v1"),
		[]byte("v2"),
		[]byte("v3"),
		[]byte("v4"),
	}
	generatedKey := "generated_key"
	dtUpdateResp := &rpbRiakDT.DtUpdateResp{
		SetValue: setValue,
		Key:      []byte(generatedKey),
	}

	builder := NewUpdateSetCommandBuilder().
		WithBucketType("sets").
		WithBucket("bucket").
		WithKey("key")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}

	cmd.onSuccess(dtUpdateResp)

	if uc, ok := cmd.(*UpdateSetCommand); ok {
		rsp := uc.Response
		for i := 1; i <= 4; i++ {
			sitem := fmt.Sprintf("v%d", i)
			if expected, actual := true, sliceIncludes(rsp.SetValue, []byte(sitem)); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		}
		if expected, actual := "generated_key", rsp.GeneratedKey; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *UpdateSetCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestValidationOfUpdateSetViaBuilder(t *testing.T) {
	// validate that Bucket is required
	builder := NewUpdateSetCommandBuilder()
	_, err := builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrBucketRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	// validate that Key is NOT required
	builder = NewUpdateSetCommandBuilder()
	builder.WithBucket("bucket_name")
	_, err = builder.Build()
	if err != nil {
		t.Fatal("expected nil err")
	}
}

// FetchSet
// DtFetchReq
// DtFetchResp

func TestBuildDtFetchReqCorrectlyViaFetchSetCommandBuilder(t *testing.T) {
	builder := NewFetchSetCommandBuilder().
		WithBucketType("sets").
		WithBucket("bucket").
		WithKey("key").
		WithR(1).
		WithPr(2).
		WithNotFoundOk(true).
		WithBasicQuorum(true).
		WithTimeout(time.Second * 20)
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}
	if req, ok := protobuf.(*rpbRiakDT.DtFetchReq); ok {
		if expected, actual := "sets", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "bucket", string(req.GetBucket()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "key", string(req.GetKey()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(1), req.GetR(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(2), req.GetPr(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetNotfoundOk(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetBasicQuorum(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		validateTimeout(t, time.Second*20, req.GetTimeout())
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakDT.DtFetchReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestFetchSetParsesDtFetchRespCorrectly(t *testing.T) {
	dtValue := &rpbRiakDT.DtValue{
		SetValue: [][]byte{
			[]byte("v1"),
			[]byte("v2"),
			[]byte("v3"),
			[]byte("v4"),
		},
	}
	dtFetchResp := &rpbRiakDT.DtFetchResp{
		Value: dtValue,
	}
	builder := NewFetchSetCommandBuilder().
		WithBucketType("sets").
		WithBucket("bucket").
		WithKey("key")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}

	cmd.onSuccess(dtFetchResp)

	if fc, ok := cmd.(*FetchSetCommand); ok {
		rsp := fc.Response
		for i := 1; i <= 4; i++ {
			sitem := fmt.Sprintf("v%d", i)
			if expected, actual := true, sliceIncludes(rsp.SetValue, []byte(sitem)); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *FetchSetCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestFetchSetParsesDtFetchRespWithoutValueCorrectly(t *testing.T) {
	builder := NewFetchSetCommandBuilder().
		WithBucketType("counters").
		WithBucket("myBucket").
		WithKey("counter_1")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}

	dtFetchResp := &rpbRiakDT.DtFetchResp{}
	cmd.onSuccess(dtFetchResp)

	if uc, ok := cmd.(*FetchSetCommand); ok {
		if expected, actual := true, uc.Response.IsNotFound; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *FetchSetCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestValidationOfFetchSetViaBuilder(t *testing.T) {
	// validate that Bucket is required
	builder := NewFetchSetCommandBuilder()
	_, err := builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrBucketRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	// validate that Key is required
	builder = NewFetchSetCommandBuilder()
	builder.WithBucket("bucket_name")
	_, err = builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrKeyRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

// UpdateMap
// DtUpdateReq
// DtUpdateResp

func TestBuildDtUpdateReqCorrectlyViaUpdateMapCommandBuilder(t *testing.T) {
	mapOp := &MapOperation{}
	mapOp.IncrementCounter("counter_1", 50).
		RemoveCounter("counter_2").
		AddToSet("set_1", []byte("set_value_1")).
		RemoveFromSet("set_2", []byte("set_value_2")).
		RemoveSet("set_3").
		SetRegister("register_1", []byte("register_value_1")).
		RemoveRegister("register_2").
		SetFlag("flag_1", true).
		RemoveFlag("flag_2").
		RemoveMap("map_3")

	mapOp.Map("map_2").
		IncrementCounter("counter_1", 50).
		RemoveCounter("counter_2").
		AddToSet("set_1", []byte("set_value_1")).
		RemoveFromSet("set_2", []byte("set_value_2")).
		RemoveSet("set_3").
		SetRegister("register_1", []byte("register_value_1")).
		RemoveRegister("register_2").
		SetFlag("flag_1", true).
		RemoveFlag("flag_2").
		RemoveMap("map_3")

	builder := NewUpdateMapCommandBuilder().
		WithBucketType("maps").
		WithBucket("bucket").
		WithKey("key").
		WithContext(crdtContextBytes).
		WithMapOperation(mapOp).
		WithW(3).
		WithPw(1).
		WithDw(2).
		WithReturnBody(true).
		WithTimeout(time.Second * 20)
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}
	if req, ok := protobuf.(*rpbRiakDT.DtUpdateReq); ok {
		if expected, actual := "maps", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "bucket", string(req.GetBucket()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "key", string(req.GetKey()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := 0, bytes.Compare(crdtContextBytes, req.GetContext()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		validateTimeout(t, time.Second*20, req.GetTimeout())

		mapOp := req.Op.MapOp

		verifyRemoves := func(removes []*rpbRiakDT.MapField) {
			if expected, actual := 5, len(removes); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			counterRemoved := false
			setRemoved := false
			registerRemoved := false
			flagRemoved := false
			mapRemoved := false
			for _, remove := range removes {
				switch remove.GetType() {
				case rpbRiakDT.MapField_COUNTER:
					if expected, actual := "counter_2", string(remove.Name); expected != actual {
						t.Errorf("expected %v, got %v", expected, actual)
					}
					counterRemoved = true
				case rpbRiakDT.MapField_SET:
					if expected, actual := "set_3", string(remove.Name); expected != actual {
						t.Errorf("expected %v, got %v", expected, actual)
					}
					setRemoved = true
				case rpbRiakDT.MapField_MAP:
					if expected, actual := "map_3", string(remove.Name); expected != actual {
						t.Errorf("expected %v, got %v", expected, actual)
					}
					mapRemoved = true
				case rpbRiakDT.MapField_REGISTER:
					if expected, actual := "register_2", string(remove.Name); expected != actual {
						t.Errorf("expected %v, got %v", expected, actual)
					}
					registerRemoved = true
				case rpbRiakDT.MapField_FLAG:
					if expected, actual := "flag_2", string(remove.Name); expected != actual {
						t.Errorf("expected %v, got %v", expected, actual)
					}
					flagRemoved = true
				}
			}
			if expected, actual := true, counterRemoved; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := true, setRemoved; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := true, registerRemoved; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := true, flagRemoved; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := true, mapRemoved; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		}

		verifyUpdates := func(updates []*rpbRiakDT.MapUpdate, expectMapUpdate bool) *rpbRiakDT.MapUpdate {
			counterIncremented := false
			setAddedTo := false
			setRemovedFrom := false
			registerSet := false
			flagSet := false
			mapAdded := false
			var mapUpdate *rpbRiakDT.MapUpdate
			for _, update := range updates {
				field := update.GetField()
				switch field.GetType() {
				case rpbRiakDT.MapField_COUNTER:
					if expected, actual := "counter_1", string(field.GetName()); expected != actual {
						t.Errorf("expected %v, got %v", expected, actual)
					}
					if expected, actual := int64(50), update.CounterOp.GetIncrement(); expected != actual {
						t.Errorf("expected %v, got %v", expected, actual)
					}
					counterIncremented = true
				case rpbRiakDT.MapField_SET:
					if len(update.SetOp.Adds) > 0 {
						if expected, actual := "set_1", string(field.GetName()); expected != actual {
							t.Errorf("expected %v, got %v", expected, actual)
						}
						if expected, actual := "set_value_1", string(update.SetOp.Adds[0]); expected != actual {
							t.Errorf("expected %v, got %v", expected, actual)
						}
						setAddedTo = true

					} else {
						if expected, actual := "set_2", string(field.GetName()); expected != actual {
							t.Errorf("expected %v, got %v", expected, actual)
						}
						if expected, actual := "set_value_2", string(update.SetOp.Removes[0]); expected != actual {
							t.Errorf("expected %v, got %v", expected, actual)
						}
						setRemovedFrom = true
					}
				case rpbRiakDT.MapField_MAP:
					if expectMapUpdate {
						if expected, actual := "map_2", string(field.GetName()); expected != actual {
							t.Errorf("expected %v, got %v", expected, actual)
						}
						mapAdded = true
						mapUpdate = update
					}
				case rpbRiakDT.MapField_REGISTER:
					if expected, actual := "register_1", string(field.GetName()); expected != actual {
						t.Errorf("expected %v, got %v", expected, actual)
					}
					if expected, actual := "register_value_1", string(update.RegisterOp); expected != actual {
						t.Errorf("expected %v, got %v", expected, actual)
					}
					registerSet = true
				case rpbRiakDT.MapField_FLAG:
					if expected, actual := "flag_1", string(field.GetName()); expected != actual {
						t.Errorf("expected %v, got %v", expected, actual)
					}
					if expected, actual := rpbRiakDT.MapUpdate_ENABLE, update.GetFlagOp(); expected != actual {
						t.Errorf("expected %v, got %v", expected, actual)
					}
					flagSet = true
				}
			}

			if expected, actual := true, counterIncremented; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := true, setAddedTo; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := true, setRemovedFrom; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := true, registerSet; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := true, flagSet; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expectMapUpdate {
				if expected, actual := true, mapAdded; expected != actual {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			} else {
				if expected, actual := false, mapAdded; expected != actual {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			}

			return mapUpdate
		}

		verifyRemoves(mapOp.GetRemoves())
		innerMapUpdate := verifyUpdates(mapOp.GetUpdates(), true)
		verifyRemoves(innerMapUpdate.MapOp.GetRemoves())
		verifyUpdates(innerMapUpdate.MapOp.GetUpdates(), false)

	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakDT.DtUpdateReq", ok, reflect.TypeOf(protobuf))
	}
}

/*
TODO
func TestUpdateMapParsesDtUpdateRespCorrectly(t *testing.T) {
	setValue := [][]byte{
		[]byte("v1"),
		[]byte("v2"),
		[]byte("v3"),
		[]byte("v4"),
	}
	generatedKey := "generated_key"
	dtUpdateResp := &rpbRiakDT.DtUpdateResp{
		SetValue: setValue,
		Key:      []byte(generatedKey),
	}

	builder := NewUpdateMapCommandBuilder().
		WithBucketType("sets").
		WithBucket("bucket").
		WithKey("key")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}

	cmd.onSuccess(dtUpdateResp)

	if uc, ok := cmd.(*UpdateMapCommand); ok {
		rsp := uc.Response
		for i := 1; i <= 4; i++ {
			sitem := fmt.Sprintf("v%d", i)
			if expected, actual := true, sliceIncludes(rsp.SetValue, []byte(sitem)); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		}
		if expected, actual := "generated_key", rsp.GeneratedKey; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *UpdateMapCommand", ok, reflect.TypeOf(cmd))
	}
}
*/

func TestValidationOfUpdateMapViaBuilder(t *testing.T) {
	// validate that Bucket is required
	builder := NewUpdateMapCommandBuilder()
	_, err := builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrBucketRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	// validate that Key is NOT required
	builder = NewUpdateMapCommandBuilder()
	builder.WithBucket("bucket_name")
	_, err = builder.Build()
	if err != nil {
		t.Fatal("expected nil err")
	}
}