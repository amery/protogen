package protogen

// FileData contains the general fields for rendering code files
type FileData struct {
	src *File

	ProtoCGenCmd string
}

// File returns the source file associated to this data
func (fd *FileData) File() *File {
	if fd != nil && fd.src != nil {
		return fd.src
	}
	return nil
}

// ProtoFile indicates the path of the .proto being generated
func (fd *FileData) ProtoFile() string {
	return fd.src.Name()
}

// ProtoFileBase indicates the path of the .proto with the .proto
// suffix removed.
func (fd *FileData) ProtoFileBase() string {
	return fd.src.Base()
}

// ProtoParameters indicates the parameters passed by protoc
func (fd *FileData) ProtoParameters() string {
	return optional(fd.src.Request().Parameter, "")
}

// NewFileData creates [FileData] for the given [File]
func (gen *Plugin) NewFileData(src *File) (FileData, error) {
	fd := FileData{
		src: src,

		ProtoCGenCmd: gen.options.Name,
	}

	return fd, nil
}
