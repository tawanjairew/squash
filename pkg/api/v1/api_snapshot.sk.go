// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"fmt"

	"github.com/solo-io/solo-kit/pkg/utils/hashutils"
	"go.uber.org/zap"
)

type ApiSnapshot struct {
	Debugattachments DebugattachmentsByNamespace
}

func (s ApiSnapshot) Clone() ApiSnapshot {
	return ApiSnapshot{
		Debugattachments: s.Debugattachments.Clone(),
	}
}

func (s ApiSnapshot) Hash() uint64 {
	return hashutils.HashAll(
		s.hashDebugattachments(),
	)
}

func (s ApiSnapshot) hashDebugattachments() uint64 {
	return hashutils.HashAll(s.Debugattachments.List().AsInterfaces()...)
}

func (s ApiSnapshot) HashFields() []zap.Field {
	var fields []zap.Field
	fields = append(fields, zap.Uint64("debugattachments", s.hashDebugattachments()))

	return append(fields, zap.Uint64("snapshotHash", s.Hash()))
}

type ApiSnapshotStringer struct {
	Version          uint64
	Debugattachments []string
}

func (ss ApiSnapshotStringer) String() string {
	s := fmt.Sprintf("ApiSnapshot %v\n", ss.Version)

	s += fmt.Sprintf("  Debugattachments %v\n", len(ss.Debugattachments))
	for _, name := range ss.Debugattachments {
		s += fmt.Sprintf("    %v\n", name)
	}

	return s
}

func (s ApiSnapshot) Stringer() ApiSnapshotStringer {
	return ApiSnapshotStringer{
		Version:          s.Hash(),
		Debugattachments: s.Debugattachments.List().NamespacesDotNames(),
	}
}
