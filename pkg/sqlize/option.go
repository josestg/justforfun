package sqlize

// WithPrinter knows how to set Migrator printer.
func WithPrinter(printer Printer) Option {
	return func(m *Migrator) {
		m.printer = printer
	}
}

// WithTemplating know how to set Migrator templating engine.
func WithTemplating(t TemplateReader) Option {
	return func(m *Migrator) {
		m.templating = t
	}
}

// WithVersioning know how to set Migrator version generator.
func WithVersioning(fn VersionFunc) Option {
	return func(m *Migrator) {
		m.versioning = fn
	}
}
