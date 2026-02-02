package redis

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type tUsers struct {
	UserId     int       `json:"user_id" db:"user_id,primary"`
	UserName   string    `json:"user_name" db:"user_name"`
	Status     int       `json:"status" db:"status"`
	Timezone   string    `json:"timezone" db:"timezone"`
	Lang       string    `json:"lang" db:"lang"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}

type tOrders struct {
	OrderId    int          `json:"order_id" db:"order_id"`
	Body       string       `json:"body" db:"body"`
	ExtAttr    OrderExtAttr `json:"ext_attr" db:"ext_attr"` // 扩展属性存放
	CreateTime time.Time    `json:"create_time" db:"create_time"`
	UpdateTime time.Time    `json:"update_time" db:"update_time"`
}

type OrderExtAttr struct {
	NotifyUrl string `json:"notify_url,omitempty"` // 回调通知url
}

var emptyJSON = json.RawMessage("{}")

func JsonObject(value any) (json.RawMessage, error) {
	var source []byte
	switch t := value.(type) {
	case string:
		source = []byte(t)
	case []byte:
		source = t
	case nil:
		source = emptyJSON
	default:
		return nil, errors.New("incompatible type for json.RawMessage")
	}

	if len(source) == 0 {
		source = emptyJSON
	}

	return source, nil
}

// Scan implements the Scanner interface.
func (n *OrderExtAttr) Scan(value any) error {
	b, err := JsonObject(value)
	if err != nil {
		return err
	}
	// 忽略非对象形式的数据脏数据，比如：数组
	if len(b) > 0 && []byte(b)[0] == '[' {
		b = json.RawMessage("{}")
	}
	return json.Unmarshal(b, n)
}

// Value implements the driver Valuer interface.
func (n OrderExtAttr) Value() (driver.Value, error) {
	return json.Marshal(n)
}

func (n *OrderExtAttr) UnmarshalBinary(data []byte) error {
	// convert data to yours, let's assume its json data
	return json.Unmarshal(data, n)
}

func (n OrderExtAttr) MarshalBinary() ([]byte, error) {
	return json.Marshal(n)
}

var m = map[string]string{
	"user_id":     "4",
	"user_name":   "test",
	"status":      "1",
	"timezone":    "Asia/Shanghai",
	"lang":        "zh-CN",
	"create_time": "2015-03-18 18:20:28",
	"update_time": "2017-09-20 10:29:59",
}

func BenchmarkScanStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		user := &tUsers{}
		_ = scanStruct(m, user)
	}
}

func TestScanStructOrder(t *testing.T) {
	order := &tOrders{}

	m := map[string]string{
		"order_id":    "123123",
		"body":        "test",
		"ext_attr":    `{"notify_url":"https://v.com"}`,
		"create_time": "2015-03-18 18:20:28",
		"update_time": "2017-09-20 10:29:59",
	}

	err := scanStruct(m, order)
	assert.NoError(t, err)
	assert.Equal(t, 123123, order.OrderId)
	assert.Equal(t, "2015-03-18 18:20:28", order.CreateTime.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "https://v.com", order.ExtAttr.NotifyUrl)
}

func TestScanStruct(t *testing.T) {
	t.Run("tag db", func(t *testing.T) {
		user := &tUsers{}
		err := scanStruct(m, user)

		if err != nil {
			t.Error(err)
		}

		if user.UserId != 4 {
			t.Error("Parse user_id error:", user.UserId)
		}

		if user.CreateTime.Format("2006-01-02 15:04:05") != "2015-03-18 18:20:28" {
			t.Error("Parse create_time error:", user.CreateTime)
		}
	})

	t.Run("tag json", func(t *testing.T) {
		user := &tUsers{}
		err := scanStructWithTag(m, user, "json")

		if err != nil {
			t.Error(err)
		}

		if user.UserId != 4 {
			t.Error("Parse user_id error:", user.UserId)
		}

		if user.CreateTime.Format("2006-01-02 15:04:05") != "2015-03-18 18:20:28" {
			t.Error("Parse create_time error:", user.CreateTime)
		}
	})

	t.Run("tag emtpy", func(t *testing.T) {
		var m = map[string]string{
			"UserId":     "4",
			"UserName":   "test",
			"Status":     "1",
			"Timezone":   "Asia/Shanghai",
			"Lang":       "zh-CN",
			"CreateTime": "2015-03-18 18:20:28",
			"UpdateTime": "2017-09-20 10:29:59",
		}

		user := &tUsers{}
		err := scanStructWithTag(m, user, "")

		if err != nil {
			t.Error(err)
		}

		if user.UserId != 4 {
			t.Error("Parse user_id error:", user.UserId)
		}

		if user.CreateTime.Format("2006-01-02 15:04:05") != "2015-03-18 18:20:28" {
			t.Error("Parse create_time error:", user.CreateTime)
		}
	})
}

func Test_structToMapInterface(t *testing.T) {
	user := &tUsers{
		UserId:     1000,
		CreateTime: time.Date(2019, 11, 02, 15, 04, 05, 0, time.UTC),
	}
	m := structToMap(user, ScanTagName)
	var now2 string

	if m2, has := m["create_time"]; has {
		if m3, ok := m2.(string); ok {
			now2 = m3
		}
	}

	if now2 != "2019-11-02 15:04:05" {
		t.Error("struct to map time parse error")
	}

	assert.Equal(t, m["user_id"], user.UserId)
}

func Test_structToMap(t *testing.T) {
	type args struct {
		Id   int
		Name string
	}

	s := &args{
		Id:   1,
		Name: "test",
	}

	m := structToMap(s, "")

	assert.Equal(t, m["Id"], s.Id)
	assert.Equal(t, m["Name"], s.Name)
}
