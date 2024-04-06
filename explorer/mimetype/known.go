package mimetype

// Any is a mime type that matches any mime type.
var Any = MimeType{Extension: "*", Type: "*", Subtype: "*"}

var (
	// TextAny is a mime type that matches any text mime type.
	TextAny = MimeType{Extension: "*", Type: "text", Subtype: "*"}

	// TextPlain is the mime type for plain text.
	TextPlain = MimeType{Extension: "txt", Type: "text", Subtype: "plain"}

	// TextHTML is the mime type for HTML.
	TextHTML = MimeType{Extension: "html", Type: "text", Subtype: "html"}

	// TextCSS is the mime type for CSS.
	TextCSS = MimeType{Extension: "css", Type: "text", Subtype: "css"}

	// TextXML is the mime type for XML.
	TextXML = MimeType{Extension: "xml", Type: "text", Subtype: "xml"}

	// TextYAML is the mime type for YAML.
	TextYAML = MimeType{Extension: "yaml", Type: "text", Subtype: "yaml"}
)

var (
	// ImageAny is the mime type for any image.
	ImageAny = MimeType{Extension: "*", Type: "image", Subtype: "*"}

	// ImageGIF is the mime type for GIF.
	ImageGIF = MimeType{Extension: "gif", Type: "image", Subtype: "gif"}

	// ImageJPEG is the mime type for JPEG.
	ImageJPEG = MimeType{Extension: "jpg", Type: "image", Subtype: "jpeg"}

	// ImagePNG is the mime type for PNG.
	ImagePNG = MimeType{Extension: "png", Type: "image", Subtype: "png"}

	// ImageWebP is the mime type for WebP.
	ImageWebP = MimeType{Extension: "webp", Type: "image", Subtype: "webp"}

	// ImageBMP is the mime type for BMP.
	ImageBMP = MimeType{Extension: "bmp", Type: "image", Subtype: "bmp"}

	// ImageTIFF is the mime type for TIFF.
	ImageTIFF = MimeType{Extension: "tiff", Type: "image", Subtype: "tiff"}

	// ImageSVG is the mime type for SVG.
	ImageSVG = MimeType{Extension: "svg", Type: "image", Subtype: "svg"}

	// ImageICO is the mime type for ICO.
	ImageICO = MimeType{Extension: "ico", Type: "image", Subtype: "ico"}

	// ImageAPNG is the mime type for APNG.
	ImageAPNG = MimeType{Extension: "apng", Type: "image", Subtype: "apng"}

	// ImageAVIF is the mime type for AVIF.
	ImageAVIF = MimeType{Extension: "avif", Type: "image", Subtype: "avif"}

	// ImageJXL is the mime type for JXL.
	ImageJXL = MimeType{Extension: "jxl", Type: "image", Subtype: "jxl"}

	// ImageHEIF is the mime type for HEIF.
	ImageHEIF = MimeType{Extension: "heif", Type: "image", Subtype: "heif"}

	// ImageHEIC is the mime type for HEIC.
	ImageHEIC = MimeType{Extension: "heic", Type: "image", Subtype: "heic"}
)

var (
	// BinaryAny is the mime type for any binary.
	BinaryAny = MimeType{Extension: "*", Type: "application", Subtype: "*"}

	// BinaryZip is the mime type for ZIP.
	BinaryZip = MimeType{Extension: "zip", Type: "application", Subtype: "zip"}
)

var (
	// AudioAny is the mime type for any audio.
	AudioAny = MimeType{Extension: "*", Type: "audio", Subtype: "*"}

	// AudioMP3 is the mime type for MP3.
	AudioMP3 = MimeType{Extension: "mp3", Type: "audio", Subtype: "mp3"}

	// AudioOGG is the mime type for OGG.
	AudioOGG = MimeType{Extension: "ogg", Type: "audio", Subtype: "ogg"}

	// AudioWAV is the mime type for WAV.
	AudioWAV = MimeType{Extension: "wav", Type: "audio", Subtype: "wav"}

	// AudioFLAC is the mime type for FLAC.
	AudioFLAC = MimeType{Extension: "flac", Type: "audio", Subtype: "flac"}

	// AudioMIDI is the mime type for MIDI.
	AudioMIDI = MimeType{Extension: "midi", Type: "audio", Subtype: "midi"}
)

var (
	// VideoAny is the mime type for any video.
	VideoAny = MimeType{Extension: "*", Type: "video", Subtype: "*"}

	// VideoMP4 is the mime type for MP4.
	VideoMP4 = MimeType{Extension: "mp4", Type: "video", Subtype: "mp4"}

	// VideoOGG is the mime type for OGG.
	VideoOGG = MimeType{Extension: "ogg", Type: "video", Subtype: "ogg"}

	// VideoWebM is the mime type for WebM.
	VideoWebM = MimeType{Extension: "webm", Type: "video", Subtype: "webm"}

	// VideoAVI is the mime type for AVI.
	VideoAVI = MimeType{Extension: "avi", Type: "video", Subtype: "avi"}

	// VideoWMV is the mime type for WMV.
	VideoWMV = MimeType{Extension: "wmv", Type: "video", Subtype: "wmv"}

	// VideoMOV is the mime type for MOV.
	VideoMOV = MimeType{Extension: "mov", Type: "video", Subtype: "mov"}

	// VideoFLV is the mime type for FLV.
	VideoFLV = MimeType{Extension: "flv", Type: "video", Subtype: "flv"}

	// VideoMKV is the mime type for MKV.
	VideoMKV = MimeType{Extension: "mkv", Type: "video", Subtype: "mkv"}

	// VideoMPEG is the mime type for MPEG.
	VideoMPEG = MimeType{Extension: "mpeg", Type: "video", Subtype: "mpeg"}
)
