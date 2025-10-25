package logger

type DummyLogger struct{}

func (l *DummyLogger) With() LogEntry  { return &DummyLogEntry{} }
func (l *DummyLogger) Trace() LogEntry { return &DummyLogEntry{} }
func (l *DummyLogger) Debug() LogEntry { return &DummyLogEntry{} }
func (l *DummyLogger) Info() LogEntry  { return &DummyLogEntry{} }
func (l *DummyLogger) Warn() LogEntry  { return &DummyLogEntry{} }
func (l *DummyLogger) Error() LogEntry { return &DummyLogEntry{} }
func (l *DummyLogger) Fatal() LogEntry { return &DummyLogEntry{} }
func (l *DummyLogger) Panic() LogEntry { return &DummyLogEntry{} }

type DummyLogEntry struct{}

func (e *DummyLogEntry) Str(key, value string) LogEntry             { return e }
func (e *DummyLogEntry) Bool(key string, value bool) LogEntry       { return e }
func (e *DummyLogEntry) Int(key string, value int) LogEntry         { return e }
func (e *DummyLogEntry) Float64(key string, value float64) LogEntry { return e }
func (e *DummyLogEntry) Any(key string, value any) LogEntry         { return e }
func (e *DummyLogEntry) Err(err error) LogEntry                     { return e }
func (e *DummyLogEntry) Msg(msg string)                             {}
func (e *DummyLogEntry) Msgf(format string, v ...any)               {}
func (e *DummyLogEntry) Send()                                      {}
func (e *DummyLogEntry) Logger() Logger                             { return &DummyLogger{} }
