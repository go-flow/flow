package flow

import "testing"

func TestConfig(t *testing.T) {
	type configType int
	const (
		stringVal configType = 1 << iota
		InterfaceVal
		IntVal
		Int32Val
		Int64Val
		BoolVal
	)
	//create config data
	cfg := Config{}
	cfg["key1"] = "value1"
	cfg["key2"] = 2
	cfg["key3"] = Int64Val
	cfg["key4"] = int32(2)
	cfg["key5"] = int64(2)
	cfg["key6"] = true

	// create config object
	//cfg := Config{data}

	// create test cases
	tt := []struct {
		Name        string      // test case name
		Type        configType  // test case type
		Key         string      // key to be used
		DefaultVal  interface{} // default value
		ExpectedVal interface{} // expected value
	}{
		{Name: "Get Interface", Type: InterfaceVal, Key: "key3", DefaultVal: nil, ExpectedVal: Int64Val},
		{Name: "Get Default Interface default", Type: InterfaceVal, Key: "key3", DefaultVal: Int32Val, ExpectedVal: Int64Val},
		{Name: "Get Default Interface no key", Type: InterfaceVal, Key: "key34", DefaultVal: nil, ExpectedVal: nil},
		{Name: "Get Default Interface no key default", Type: InterfaceVal, Key: "key44", DefaultVal: Int64Val, ExpectedVal: Int64Val},

		{Name: "Get String", Type: stringVal, Key: "key1", DefaultVal: nil, ExpectedVal: "value1"},
		{Name: "Get String default", Type: stringVal, Key: "key1", DefaultVal: "123", ExpectedVal: "value1"},
		{Name: "Get Default String no key", Type: stringVal, Key: "ket1", DefaultVal: nil, ExpectedVal: ""},
		{Name: "Get Default Interface no key default", Type: stringVal, Key: "key44", DefaultVal: "value44", ExpectedVal: "value44"},

		{Name: "Get Int", Type: IntVal, Key: "key2", DefaultVal: nil, ExpectedVal: 2},
		{Name: "Get Int default", Type: IntVal, Key: "key2", DefaultVal: 1, ExpectedVal: 2},
		{Name: "Get Int wrong type", Type: IntVal, Key: "key1", DefaultVal: nil, ExpectedVal: 0},
		{Name: "Get Int wrong type default", Type: IntVal, Key: "key1", DefaultVal: 12, ExpectedVal: 12},
		{Name: "Get Default Int no key", Type: IntVal, Key: "key44", DefaultVal: nil, ExpectedVal: 0},
		{Name: "Get Default Int no key default", Type: IntVal, Key: "key44", DefaultVal: 12, ExpectedVal: 12},

		{Name: "Get Int32", Type: Int32Val, Key: "key4", DefaultVal: nil, ExpectedVal: int32(2)},
		{Name: "Get Int32 default", Type: Int32Val, Key: "key4", DefaultVal: int32(1), ExpectedVal: int32(2)},
		{Name: "Get Int32 wrong type", Type: Int32Val, Key: "key1", DefaultVal: nil, ExpectedVal: int32(0)},
		{Name: "Get Int32 wrong type default", Type: Int32Val, Key: "key1", DefaultVal: int32(12), ExpectedVal: int32(12)},
		{Name: "Get Default Int32 no key", Type: Int32Val, Key: "key44", DefaultVal: nil, ExpectedVal: int32(0)},
		{Name: "Get Default Int32 no key default", Type: Int32Val, Key: "key44", DefaultVal: int32(12), ExpectedVal: int32(12)},

		{Name: "Get Int64", Type: Int64Val, Key: "key5", DefaultVal: nil, ExpectedVal: int64(2)},
		{Name: "Get Int64 default", Type: Int64Val, Key: "key5", DefaultVal: int64(1), ExpectedVal: int64(2)},
		{Name: "Get Int64 wrong type", Type: Int64Val, Key: "key1", DefaultVal: nil, ExpectedVal: int64(0)},
		{Name: "Get Int64 wrong type default", Type: Int64Val, Key: "key1", DefaultVal: int64(12), ExpectedVal: int64(12)},
		{Name: "Get Default Int64 no key", Type: Int64Val, Key: "key44", DefaultVal: nil, ExpectedVal: int64(0)},
		{Name: "Get Default Int64 no key default", Type: Int64Val, Key: "key44", DefaultVal: int64(12), ExpectedVal: int64(12)},

		{Name: "Get Bool", Type: BoolVal, Key: "key6", DefaultVal: nil, ExpectedVal: true},
		{Name: "Get Bool default", Type: BoolVal, Key: "key6", DefaultVal: false, ExpectedVal: true},
		{Name: "Get Bool wrong type", Type: BoolVal, Key: "key1", DefaultVal: nil, ExpectedVal: false},
		{Name: "Get Bool wrong type default", Type: BoolVal, Key: "key1", DefaultVal: true, ExpectedVal: true},
		{Name: "Get Default Bool no key", Type: BoolVal, Key: "key44", DefaultVal: nil, ExpectedVal: false},
		{Name: "Get Default Bool no key default", Type: BoolVal, Key: "key44", DefaultVal: true, ExpectedVal: true},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			switch tc.Type {
			case InterfaceVal:
				if tc.DefaultVal == nil {
					r := cfg.Get(tc.Key)
					if r != tc.ExpectedVal {
						t.Errorf("%s - error expected %v, got %v", tc.Name, tc.ExpectedVal, r)
					}
				} else {
					r := cfg.GetDefault(tc.Key, tc.DefaultVal)
					if r != tc.ExpectedVal {
						t.Errorf("%s - error expected %v, got %v", tc.Name, tc.ExpectedVal, r)
					}
				}
			case stringVal:
				if tc.DefaultVal == nil {
					r := cfg.String(tc.Key)
					if r != tc.ExpectedVal {
						t.Errorf("%s - error expected %v, got %v", tc.Name, tc.ExpectedVal, r)
					}
				} else {
					r := cfg.StringDefault(tc.Key, tc.DefaultVal.(string))
					if r != tc.ExpectedVal {
						t.Errorf("%s - error expected %v, got %v", tc.Name, tc.ExpectedVal, r)
					}
				}
			case IntVal:
				if tc.DefaultVal == nil {
					r := cfg.Int(tc.Key)
					if r != tc.ExpectedVal {
						t.Errorf("%s - error expected %v, got %v", tc.Name, tc.ExpectedVal, r)
					}
				} else {
					r := cfg.IntDefault(tc.Key, tc.DefaultVal.(int))
					if r != tc.ExpectedVal {
						t.Errorf("%s - error expected %v, got %v", tc.Name, tc.ExpectedVal, r)
					}
				}
			case Int32Val:
				if tc.DefaultVal == nil {
					r := cfg.Int32(tc.Key)
					if r != tc.ExpectedVal {
						t.Errorf("%s - error expected %v, got %v", tc.Name, tc.ExpectedVal, r)
					}
				} else {
					r := cfg.Int32Default(tc.Key, tc.DefaultVal.(int32))
					if r != tc.ExpectedVal {
						t.Errorf("%s - error expected %v, got %v", tc.Name, tc.ExpectedVal, r)
					}
				}
			case Int64Val:
				if tc.DefaultVal == nil {
					r := cfg.Int64(tc.Key)
					if r != tc.ExpectedVal {
						t.Errorf("%s - error expected %v, got %v", tc.Name, tc.ExpectedVal, r)
					}
				} else {
					r := cfg.Int64Default(tc.Key, tc.DefaultVal.(int64))
					if r != tc.ExpectedVal {
						t.Errorf("%s - error expected %v, got %v", tc.Name, tc.ExpectedVal, r)
					}
				}
			case BoolVal:
				if tc.DefaultVal == nil {
					r := cfg.Bool(tc.Key)
					if r != tc.ExpectedVal {
						t.Errorf("%s - error expected %v, got %v", tc.Name, tc.ExpectedVal, r)
					}
				} else {
					r := cfg.BoolDefault(tc.Key, tc.DefaultVal.(bool))
					if r != tc.ExpectedVal {
						t.Errorf("%s - error expected %v, got %v", tc.Name, tc.ExpectedVal, r)
					}
				}
			}

		})
	}
}
