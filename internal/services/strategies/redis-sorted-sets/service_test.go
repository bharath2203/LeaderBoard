package redis_sorted_sets

import (
	"TopKScores/api"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/go-redis/redismock/v7"
	"reflect"
	"testing"
)

func Test_redisSortedSetsImpl_validateRecordCount(t *testing.T) {
	type fields struct {
		redisClient        *redis.Client
		gameKey            string
		MaxNumberOfRecords int64
	}
	type args struct {
		k int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "greater than set number",
			fields: fields{
				redisClient:        &redis.Client{},
				gameKey:            "key",
				MaxNumberOfRecords: 10,
			},
			args: args{
				k: 12,
			},
			wantErr: true,
		},
		{
			name: "greater than set number",
			fields: fields{
				redisClient:        &redis.Client{},
				gameKey:            "key",
				MaxNumberOfRecords: 10,
			},
			args: args{
				k: 2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &redisSortedSetsImpl{
				redisClient:        tt.fields.redisClient,
				redisKey:           tt.fields.gameKey,
				maxNumberOfRecords: tt.fields.MaxNumberOfRecords,
			}
			if err := c.validateRecordCount(tt.args.k); (err != nil) != tt.wantErr {
				t.Errorf("validateRecordCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisSortedSetsImpl_AddScore(t *testing.T) {
	type fields struct {
		redisClient        *redis.Client
		gameKey            string
		MaxNumberOfRecords int64
	}
	type args struct {
		gameScore *api.GameScore
	}

	rc, mock := redismock.NewClientMock()

	mock.ExpectZAdd("test-key-1", &redis.Z{
		Score: 100,
		Member: api.GameInstance{
			GameId:   "test-id-1",
			UserName: "test-username-1",
		},
	}).SetVal(1)
	mock.ExpectZRemRangeByRank("test-key-1", 0, -2).SetVal(1)

	mock.ExpectZAdd("test-key-2", &redis.Z{
		Score: 100,
		Member: api.GameInstance{
			GameId:   "test-id-2",
			UserName: "test-username-2",
		},
	}).SetErr(fmt.Errorf("redis error"))
	mock.ExpectZRemRangeByRank("test-key-2", 0, -2).SetVal(1)

	mock.ExpectZAdd("test-key-3", &redis.Z{
		Score: 100,
		Member: api.GameInstance{
			GameId:   "test-id-3",
			UserName: "test-username-3",
		},
	}).SetVal(2)
	mock.ExpectZRemRangeByRank("test-key-3", 0, -3).SetErr(fmt.Errorf("redis error"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test-1(no err from redis)",
			fields: fields{
				redisClient:        rc,
				gameKey:            "test-key-1",
				MaxNumberOfRecords: 2,
			},
			args: args{
				gameScore: &api.GameScore{
					GameId:    "test-id-1",
					UserName:  "test-username-1",
					UserScore: 100,
				},
			},
			wantErr: false,
		},
		{
			name: "test-2(redis error on ZADD)",
			fields: fields{
				redisClient:        rc,
				gameKey:            "test-key-2",
				MaxNumberOfRecords: 2,
			},
			args: args{
				gameScore: &api.GameScore{
					GameId:    "test-id-2",
					UserName:  "test-username-2",
					UserScore: 100,
				},
			},
			wantErr: true,
		},
		{
			name: "test-2(redis error on ZREMBYRANK)",
			fields: fields{
				redisClient:        rc,
				gameKey:            "test-key-3",
				MaxNumberOfRecords: 3,
			},
			args: args{
				gameScore: &api.GameScore{
					GameId:    "test-id-3",
					UserName:  "test-username-3",
					UserScore: 100,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &redisSortedSetsImpl{
				redisClient:        tt.fields.redisClient,
				redisKey:           tt.fields.gameKey,
				maxNumberOfRecords: tt.fields.MaxNumberOfRecords,
			}
			if err := c.AddScore(tt.args.gameScore); (err != nil) != tt.wantErr {
				t.Errorf("AddScore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisSortedSetsImpl_GetTopKScores(t *testing.T) {
	testGameKey := "test-game-key"
	type fields struct {
		redisClient        *redis.Client
		gameKey            string
		maxNumberOfRecords int64
	}

	rc, mock := redismock.NewClientMock()

	mock.ExpectZRevRangeWithScores(testGameKey, 0, 1).SetVal([]redis.Z{
		{
			Score:  100,
			Member: "{\"game_id\":\"1\",\"user_name\":\"one\"}",
		},
		{
			Score:  50,
			Member: "{\"game_id\":\"2\",\"user_name\":\"two\"}",
		},
	})

	mock.ExpectZRevRangeWithScores(testGameKey, 0, 2).SetErr(fmt.Errorf("redis internal error"))

	type args struct {
		k int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*api.GameScore
		wantErr bool
	}{
		{
			name: "test 1(validation error)",
			fields: fields{
				redisClient:        rc,
				gameKey:            testGameKey,
				maxNumberOfRecords: 3,
			},
			args: args{
				k: 4,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test 2",
			fields: fields{
				redisClient:        rc,
				gameKey:            testGameKey,
				maxNumberOfRecords: 100,
			},
			args: args{
				k: 2,
			},
			want: []*api.GameScore{
				{
					GameId:    "1",
					UserName:  "one",
					UserScore: 100,
				},
				{
					GameId:    "2",
					UserName:  "two",
					UserScore: 50,
				},
			},
			wantErr: false,
		},
		{
			name: "test 3(redis error)",
			fields: fields{
				redisClient:        rc,
				gameKey:            testGameKey,
				maxNumberOfRecords: 5,
			},
			args: args{
				k: 3,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &redisSortedSetsImpl{
				redisClient:        tt.fields.redisClient,
				redisKey:           tt.fields.gameKey,
				maxNumberOfRecords: tt.fields.maxNumberOfRecords,
			}
			got, err := c.GetTopKScores(tt.args.k)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTopKScores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTopKScores() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisSortedSetsImpl_validateGameScoreObject(t *testing.T) {
	type fields struct {
		redisClient        *redis.Client
		redisKey           string
		maxNumberOfRecords int64
	}
	type args struct {
		gameScore *api.GameScore
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "object nil",
			fields: fields{
				redisClient: nil,
			},
			args: args{
				gameScore: nil,
			},
			wantErr: true,
		},
		{
			name: "game id empty",
			fields: fields{
				redisClient: nil,
			},
			args: args{
				gameScore: &api.GameScore{
					GameId:    "",
					UserName:  "test",
					UserScore: 100,
				},
			},
			wantErr: true,
		},
		{
			name: "user name empty",
			fields: fields{
				redisClient: nil,
			},
			args: args{
				gameScore: &api.GameScore{
					GameId:    "1",
					UserName:  "",
					UserScore: 100,
				},
			},
			wantErr: true,
		},
		{
			name: "game score 0",
			fields: fields{
				redisClient: nil,
			},
			args: args{
				gameScore: &api.GameScore{
					GameId:    "1",
					UserName:  "test",
					UserScore: 0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &redisSortedSetsImpl{
				redisClient:        tt.fields.redisClient,
				redisKey:           tt.fields.redisKey,
				maxNumberOfRecords: tt.fields.maxNumberOfRecords,
			}
			if err := c.validateGameScoreObject(tt.args.gameScore); (err != nil) != tt.wantErr {
				t.Errorf("validateGameScoreObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
