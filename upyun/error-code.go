package upyunBlob

var respErrorCode = map[string]string{
	"":         "unknown",
	"40000001": "need save-key",
	"40000002": "need body",
	"40000003": "need bucket name",
	"40000004": "need expiration",
	"40000005": "need file",
	"40000006": "content md5 not match",
	"40000007": "need policy",
	"40000008": "need signature",
	"40000009": "decode policy error",
	"40000010": "data too long for ext-param",
	"40000011": "chunked request bodies not supported yet",
	"40000012": "write file to fs error",
	"40000013": "need content-md5 but no body provided",
	"40000014": "need content-md5 but no content provided",
	"40000015": "missing required arguments",
	"40000016": "missing file_hash argument",
	"40000017": "missing file_blocks argument",
	"40000018": "too many file blocks",
	"40000019": "missing file_size argument",
	"40000020": "missing save_token argument",
	"40000021": "missing block_index argument",
	"40000022": "missing block_hash argument",
	"40000023": "block already exists",
	"40000024": "block not finished",
	"40000025": "file size not match",
	"40000026": "block index out of range",
	"40000027": "block size too small, at least 100KB",
	"40000028": "block size too large, at most 5MB",
	"40000029": "save-key encoding should be utf8",
	"40000030": "path encoding should be utf8",
	"40000031": "filename should be utf8",
	"40000032": "ffmpeg args error",
	"40000033": "failed to read firstchunk",
	"40000034": "client error",
	"40000035": "need purge body",
	"40000036": "uri bucket must be same as param bucket",
	"40000037": "unknown service",
	"40000038": "no boundary defined in Content-Type",
	"40010001": "not enough arguments",
	"40010002": "url not allowed",
	"40010003": "tasks number error",
	"40010029": "only 1~10 tasks supported in get",
	"40010033": "decode tasks error",
	"40011001": "invalid x-gmkerl-type",
	"40011002": "invalid x-gmkerl-value",
	"40011003": "invalid x-gmkerl-unsharp",
	"40011004": "invalid x-gmkerl-quality",
	"40011005": "invalid x-gmkerl-exif-switch",
	"40011006": "invalid x-gmkerl-gifto",
	"40011007": "invalid x-gmkerl-webp",
	"40011008": "invalid x-gmkerl-format",
	"40011009": "invalid x-gmkerl-crop",
	"40011010": "invalid x-gmkerl-rotate",
	"40011011": "invalid x-gmkerl-gaussblur",
	"40011012": "invalid x-gmkerl-compress",
	"40011013": "invalid x-gmkerl-progressive",
	"40011014": "invalid x-gmkerl-noicc",
	"40011015": "invalid x-gmkerl-watermark-switch",
	"40011016": "invalid x-gmkerl-watermark-align",
	"40011017": "invalid x-gmkerl-watermark-margin",
	"40011018": "invalid x-gmkerl-watermark-opacity",
	"40011019": "invalid x-gmkerl-watermark-animate",
	"40011020": "invalid x-gmkerl-watermark-font",
	"40011021": "invalid x-gmkerl-watermark-text",
	"40011022": "invalid x-gmkerl-watermark-color",
	"40011023": "invalid x-gmkerl-watermark-border",
	"40011024": "invalid x-gmkerl-watermark-size",
	"40011025": "invalid x-gmkerl-gradient-orientation",
	"40011026": "invalid x-gmkerl-gradient-pos",
	"40011027": "invalid x-gmkerl-gradient-startcolor",
	"40011028": "invalid x-gmkerl-gradient-stopcolor",
	"40011029": "invalid x-gmkerl-extract-func",
	"40011030": "invalid x-gmkerl-extract-color-count",
	"40011031": "invalid x-gmkerl-extract-offset",
	"40011032": "invalid x-gmkerl-extract-palette",
	"40011033": "invalid x-gmkerl-canvas",
	"40011034": "invalid x-gmkerl-canvas-color",
	"40011051": "invalid x-upyun-part-id",
	"40011052": "invalid x-upyun-multi-length",
	"40011053": "invalid x-upyun-multi-type",
	"40011054": "invalid x-upyun-part-size",
	"40011055": "invalid invalid multipart stage",
	"40011056": "invalid invalid x-upyun-multi-uuid",
	"40011057": "invalid UTF8 Key",
	"40011058": "invalid x-upyun-multi-uuid not found",
	"40011059": "file already upload",
	"40011060": "file md5 not match",
	"40011061": "part id error",
	"40011062": "part already complete",
	"40011091": "invalid app name",
	"40011092": "invalid file ttl",
	"40100001": "need date header",
	"40100002": "date offset error",
	"40100003": "unknown realm in authorization header",
	"40100004": "need authorization header",
	"40100005": "signature error",
	"40100006": "user not exists",
	"40100007": "user blocked",
	"40100008": "user blocked in this bucket",
	"40100009": "user password error",
	"40100010": "account not exist",
	"40100011": "account blocked",
	"40100012": "bucket not exist",
	"40100013": "bucket blocked",
	"40100014": "bucket removed",
	"40100015": "bucket read only",
	"40100016": "invalid date value in header",
	"40100017": "user need permission",
	"40100018": "account inactivate",
	"40100019": "account forbidden",
	"40100020": "account reject",
	"40100021": "overdue account",
	"40300001": "file name contains invalid chars (\\r\\n\\t)",
	"40300002": "file path too long",
	"40300003": "file name too long",
	"40300004": "bucket is full",
	"40300005": "directory not empty",
	"40300006": "authorization has expired",
	"40300007": "content md5 not match",
	"40300008": "file too small",
	"40300009": "file too large",
	"40300010": "file type error",
	"40300011": "has no permission to delete",
	"40300012": "need content-type",
	"40300013": "form api disabled",
	"40300014": "order must be asc or desc",
	"40300015": "path is not a file, maybe a directory",
	"40300016": "content-secret only accept [a-zA-Z0-9]",
	"40300017": "no thumb setting found for bucket",
	"40300018": "not an image",
	"40300019": "image width too small or too big",
	"40300020": "image height too small or too big",
	"40300021": "image limit exceeded",
	"40300022": "invalid image-width-range",
	"40300023": "invalid image-height-range",
	"40300024": "exceed max size",
	"40300025": "wrong content-type header",
	"40300026": "need content-length header",
	"40300027": "request body too big",
	"40300028": "request has expired",
	"40300029": "purge too much items",
	"40300030": "wrong content-length header",
	"40400001": "file or directory not found",
	"40400002": "base64 decoded err",
	"40600001": "dir not acceptable",
	"40600002": "folder already exists",
	"40800001": "read client request timeout",
	"41500001": "media type error, need content-type",
	"42900001": "too many requests",
	"42900002": "too many requests of the same uri",
	"42900003": "too many requests of the same bucket",
	"42900004": "request banned",
	"50300000": "unknown error",
	"50300001": "write hub failed",
	"50300002": "decode error",
	"50300003": "ffmpeg error",
	"50300006": "get hub error",
	"50300007": "capture file from fs error",
	"50300008": "delete hub error",
	"50300009": "get body temp file error",
	"50300010": "read body temp file error",
	"50300011": "capture file from ffmpeg error",
	"50300012": "service error in thumb",
	"50300013": "put file capture timeout",
	"50300014": "get client socket error",
	"50300015": "decode image info error",
	"50300016": "get image info error",
	"50300017": "failed to new md5",
	"50300018": "failed to new form",
	"50300020": "put file capture error",
	"50300021": "put capture broken pipe",
	"50300022": "new resource error",
	"50300023": "connect to db error",
	"50300024": "get master failed",
	"50300025": "write info error",
	"50300026": "get data error",
	"50300027": "put file to fs error",
	"50300028": "get file from fs error",
	"50300029": "unknow error",
	"50300030": "upstream closed connection",
	"50300031": "read upstream timeout",
	"50300032": "wrong data",
	"50300033": "new resource error",
	"50300034": "get info errors",
	"50300035": "decode hub body error",
	"50300036": "send to message queue error",
	"50300038": "failed to new http",
	"50300039": "failed to get client body reader",
	"50300040": "wrong config",
	"50300041": "connect to database error",
	"50300042": "delete path error",
	"50300043": "put task to message queue error",
}
