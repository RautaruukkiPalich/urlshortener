package log


/*
outputLogStruct

[2006-01-02T15:04:05-0700] [slog.Level] [msg] {
	"trace": [
			"1: op": "func1",
			"2: op": "func2",
			"3: op": "func3",
			...
	],
	"attrs": [
		"smth1": "attr1",
		"smth2": "attr2",
		"group1": {
			"groupAttr1": "smth1",
			"groupAttr2": "smth2",
		},
		...
	],
	"error": Error.Error(),
}
*/