package loggerloader

import (
	"log/slog"
)

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	switch a.Value.Kind() {
	case slog.KindAny:
		switch a.Value.Any().(type) {
		case error:
			a = slog.String(
				a.Key,
				a.Value.String(),
		)
	}}

	return a
}

// type stackFrame struct {
// 	Func   string `json:"func"`
// 	Source string `json:"source"`
// 	Line   int    `json:"line"`
// }

// func marshalStack(err error) []stackFrame {
// 	trace := xerrors.StackTrace(err)

// 	if len(trace) == 0 {
// 		return nil
// 	}

// 	frames := trace.Frames()

// 	s := make([]stackFrame, len(frames))

// 	for i, v := range frames {
// 		f := stackFrame{
// 			Source: filepath.Join(
// 				filepath.Base(filepath.Dir(v.File)),
// 				filepath.Base(v.File),
// 			),
// 			Func: filepath.Base(v.Function),
// 			Line: v.Line,
// 		}

// 		s[i] = f
// 	}

// 	return s
// }

// func fmtErr(err error) slog.Value {
// 	var groupValues []slog.Attr

// 	groupValues = append(groupValues, slog.String("msg", err.Error()))

// 	return slog.GroupValue(groupValues...)
// }

