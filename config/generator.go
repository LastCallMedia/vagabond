package config

import (
	"bufio"
	"bytes"
	"log"
	"text/template"
)

type Generator struct {
	Template     string
	TemplateName string
}

func (gen *Generator) Generate(env *Environment) (block []byte) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	t := template.Must(template.New(gen.TemplateName).Parse(gen.Template))
	err := t.Execute(writer, env)

	if err != nil {
		log.Fatal(err)
	}
	writer.Flush()

	return buf.Bytes()
}
