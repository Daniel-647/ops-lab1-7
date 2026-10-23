// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// )

// // 起止节次都包含在时段内，因此 start == end 表示合法的一节课。
// func TestConflictEndpointSinglePeriod(t *testing.T) {
// 	a := Section{ID: "A", Weekday: 1, Start: 3, End: 3, Weeks: AllWeeks}
// 	cases := []struct {
// 		name string
// 		b    Section
// 		want bool
// 	}{
// 		{"单节课与包含该节次的课程冲突", Section{ID: "B", Weekday: 1, Start: 2, End: 4, Weeks: AllWeeks}, true},
// 		{"单节课与相邻课程不冲突", Section{ID: "B", Weekday: 1, Start: 4, End: 5, Weeks: AllWeeks}, false},
// 	}
// 	for _, c := range cases {
// 		t.Run(c.name, func(t *testing.T) {
// 			body, err := json.Marshal(struct {
// 				A Section `json:"a"`
// 				B Section `json:"b"`
// 			}{A: a, B: c.b})
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			rec := httptest.NewRecorder()
// 			req := httptest.NewRequest(http.MethodPost, "/api/conflict", bytes.NewReader(body))
// 			newMux().ServeHTTP(rec, req)
// 			if rec.Code != http.StatusOK {
// 				t.Fatalf("合法的单节课时段应返回 200，实际 %d，响应 %q", rec.Code, rec.Body.String())
// 			}
// 			var got struct {
// 				Conflict bool `json:"conflict"`
// 			}
// 			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
// 				t.Fatalf("响应不是合法 JSON：%v", err)
// 			}
// 			if got.Conflict != c.want {
// 				t.Errorf("conflict = %v，期望 %v", got.Conflict, c.want)
// 			}
// 		})
// 	}
// }
