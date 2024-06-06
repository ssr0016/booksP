package pointer

import "time"

func String(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func StringDeref(ptr *string, def string) string {
	if ptr != nil {
		return *ptr
	}
	return def
}

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func LikeString(s string) string {
	if s != "" {
		return "%" + s + "%"
	}
	return ""
}

func BoolPtr(val bool) *bool {
	return &val
}

func Bool(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

func Int64Ptr(s int64) *int64 {
	return &s
}

func Int64(s *int64) int64 {
	if s != nil {
		return *s
	}
	return 0
}

func Int(s *int) int {
	if s != nil {
		return *s
	}
	return 0
}

func Float32(s *float32) float32 {
	if s != nil {
		return *s
	}
	return 0
}

func Float64(s *float64) float64 {
	if s != nil {
		return *s
	}
	return 0
}

func Float64Ptr(s float64) *float64 {
	return &s
}

func IntPtr(val int) *int {
	return &val
}

func TimeFromString(s *string) time.Time {
	if s != nil {
		t, _ := time.Parse(time.RFC3339Nano, *s)
		return t
	}
	return time.Time{}
}
