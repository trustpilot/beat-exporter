package collector

// Generated using https://mholt.github.io/json-to-go/

type ApmServer struct {
	Acm struct {
		Request struct {
			Count float64 `json:"count"`
		} `json:"request"`
		Response struct {
			Count  float64 `json:"count"`
			Errors struct {
				Closed       float64 `json:"closed"`
				Count        float64 `json:"count"`
				Decode       float64 `json:"decode"`
				Forbidden    float64 `json:"forbidden"`
				Internal     float64 `json:"internal"`
				Invalidquery float64 `json:"invalidquery"`
				Method       float64 `json:"method"`
				Notfound     float64 `json:"notfound"`
				Queue        float64 `json:"queue"`
				Ratelimit    float64 `json:"ratelimit"`
				Toolarge     float64 `json:"toolarge"`
				Unauthorized float64 `json:"unauthorized"`
				Unavailable  float64 `json:"unavailable"`
				Validate     float64 `json:"validate"`
			} `json:"errors"`
			Valid struct {
				Accepted    float64 `json:"accepted"`
				Count       float64 `json:"count"`
				Notmodified float64 `json:"notmodified"`
				Ok          float64 `json:"ok"`
			} `json:"valid"`
		} `json:"response"`
		Unset float64 `json:"unset"`
	} `json:"acm"`
	Decoder struct {
		Deflate struct {
			ContentLength float64 `json:"content-length"`
			Count         float64 `json:"count"`
		} `json:"deflate"`
		Gzip struct {
			ContentLength float64 `json:"content-length"`
			Count         float64 `json:"count"`
		} `json:"gzip"`
		MissingContentLength struct {
			Count float64 `json:"count"`
		} `json:"missing-content-length"`
		Reader struct {
			Count float64 `json:"count"`
		} `json:"reader"`
		Uncompressed struct {
			ContentLength float64 `json:"content-length"`
			Count         float64 `json:"count"`
		} `json:"uncompressed"`
	} `json:"decoder"`
	Jaeger struct {
		Grpc struct {
			Collect struct {
				Event struct {
					Dropped struct {
						Count float64 `json:"count"`
					} `json:"dropped"`
					Received struct {
						Count float64 `json:"count"`
					} `json:"received"`
				} `json:"event"`
				Request struct {
					Count float64 `json:"count"`
				} `json:"request"`
				Response struct {
					Count  float64 `json:"count"`
					Errors struct {
						Count float64 `json:"count"`
					} `json:"errors"`
					Valid struct {
						Count float64 `json:"count"`
					} `json:"valid"`
				} `json:"response"`
			} `json:"collect"`
			Sampling struct {
				Event struct {
					Dropped struct {
						Count float64 `json:"count"`
					} `json:"dropped"`
					Received struct {
						Count float64 `json:"count"`
					} `json:"received"`
				} `json:"event"`
				Request struct {
					Count float64 `json:"count"`
				} `json:"request"`
				Response struct {
					Count  float64 `json:"count"`
					Errors struct {
						Count float64 `json:"count"`
					} `json:"errors"`
					Valid struct {
						Count float64 `json:"count"`
					} `json:"valid"`
				} `json:"response"`
			} `json:"sampling"`
		} `json:"grpc"`
		HTTP struct {
			Event struct {
				Dropped struct {
					Count float64 `json:"count"`
				} `json:"dropped"`
				Received struct {
					Count float64 `json:"count"`
				} `json:"received"`
			} `json:"event"`
			Request struct {
				Count float64 `json:"count"`
			} `json:"request"`
			Response struct {
				Count  float64 `json:"count"`
				Errors struct {
					Count float64 `json:"count"`
				} `json:"errors"`
				Valid struct {
					Count float64 `json:"count"`
				} `json:"valid"`
			} `json:"response"`
		} `json:"http"`
	} `json:"jaeger"`
	Processor struct {
		Error struct {
			Frames          float64 `json:"frames"`
			Stacktraces     float64 `json:"stacktraces"`
			Transformations float64 `json:"transformations"`
		} `json:"error"`
		Metric struct {
			Transformations float64 `json:"transformations"`
		} `json:"metric"`
		Sourcemap struct {
			Counter  float64 `json:"counter"`
			Decoding struct {
				Count  float64 `json:"count"`
				Errors float64 `json:"errors"`
			} `json:"decoding"`
			Validation struct {
				Count  float64 `json:"count"`
				Errors float64 `json:"errors"`
			} `json:"validation"`
		} `json:"sourcemap"`
		Span struct {
			Frames          float64 `json:"frames"`
			Stacktraces     float64 `json:"stacktraces"`
			Transformations float64 `json:"transformations"`
		} `json:"span"`
		Stream struct {
			Accepted float64 `json:"accepted"`
			Errors   struct {
				Closed   float64 `json:"closed"`
				Invalid  float64 `json:"invalid"`
				Queue    float64 `json:"queue"`
				Server   float64 `json:"server"`
				Toolarge float64 `json:"toolarge"`
			} `json:"errors"`
		} `json:"stream"`
		Transaction struct {
			Transformations float64 `json:"transformations"`
		} `json:"transaction"`
	} `json:"processor"`
	Profile struct {
		Request struct {
			Count float64 `json:"count"`
		} `json:"request"`
		Response struct {
			Count  float64 `json:"count"`
			Errors struct {
				Closed       float64 `json:"closed"`
				Count        float64 `json:"count"`
				Decode       float64 `json:"decode"`
				Forbidden    float64 `json:"forbidden"`
				Internal     float64 `json:"internal"`
				Invalidquery float64 `json:"invalidquery"`
				Method       float64 `json:"method"`
				Notfound     float64 `json:"notfound"`
				Queue        float64 `json:"queue"`
				Ratelimit    float64 `json:"ratelimit"`
				Toolarge     float64 `json:"toolarge"`
				Unauthorized float64 `json:"unauthorized"`
				Unavailable  float64 `json:"unavailable"`
				Validate     float64 `json:"validate"`
			} `json:"errors"`
			Valid struct {
				Accepted    float64 `json:"accepted"`
				Count       float64 `json:"count"`
				Notmodified float64 `json:"notmodified"`
				Ok          float64 `json:"ok"`
			} `json:"valid"`
		} `json:"response"`
		Unset float64 `json:"unset"`
	} `json:"profile"`
	Root struct {
		Request struct {
			Count float64 `json:"count"`
		} `json:"request"`
		Response struct {
			Count  float64 `json:"count"`
			Errors struct {
				Closed       float64 `json:"closed"`
				Count        float64 `json:"count"`
				Decode       float64 `json:"decode"`
				Forbidden    float64 `json:"forbidden"`
				Internal     float64 `json:"internal"`
				Invalidquery float64 `json:"invalidquery"`
				Method       float64 `json:"method"`
				Notfound     float64 `json:"notfound"`
				Queue        float64 `json:"queue"`
				Ratelimit    float64 `json:"ratelimit"`
				Toolarge     float64 `json:"toolarge"`
				Unauthorized float64 `json:"unauthorized"`
				Unavailable  float64 `json:"unavailable"`
				Validate     float64 `json:"validate"`
			} `json:"errors"`
			Valid struct {
				Accepted    float64 `json:"accepted"`
				Count       float64 `json:"count"`
				Notmodified float64 `json:"notmodified"`
				Ok          float64 `json:"ok"`
			} `json:"valid"`
		} `json:"response"`
		Unset float64 `json:"unset"`
	} `json:"root"`
	Sampling struct {
		TransactionsDropped float64 `json:"transactions_dropped"`
	} `json:"sampling"`
	Server struct {
		Request struct {
			Count float64 `json:"count"`
		} `json:"request"`
		Response struct {
			Count  float64 `json:"count"`
			Errors struct {
				Closed       float64 `json:"closed"`
				Count        float64 `json:"count"`
				Decode       float64 `json:"decode"`
				Forbidden    float64 `json:"forbidden"`
				Internal     float64 `json:"internal"`
				Invalidquery float64 `json:"invalidquery"`
				Method       float64 `json:"method"`
				Notfound     float64 `json:"notfound"`
				Queue        float64 `json:"queue"`
				Ratelimit    float64 `json:"ratelimit"`
				Toolarge     float64 `json:"toolarge"`
				Unauthorized float64 `json:"unauthorized"`
				Unavailable  float64 `json:"unavailable"`
				Validate     float64 `json:"validate"`
			} `json:"errors"`
			Valid struct {
				Accepted    float64 `json:"accepted"`
				Count       float64 `json:"count"`
				Notmodified float64 `json:"notmodified"`
				Ok          float64 `json:"ok"`
			} `json:"valid"`
		} `json:"response"`
		Unset float64 `json:"unset"`
	} `json:"server"`
	Sourcemap struct {
		Request struct {
			Count float64 `json:"count"`
		} `json:"request"`
		Response struct {
			Count  float64 `json:"count"`
			Errors struct {
				Closed       float64 `json:"closed"`
				Count        float64 `json:"count"`
				Decode       float64 `json:"decode"`
				Forbidden    float64 `json:"forbidden"`
				Internal     float64 `json:"internal"`
				Invalidquery float64 `json:"invalidquery"`
				Method       float64 `json:"method"`
				Notfound     float64 `json:"notfound"`
				Queue        float64 `json:"queue"`
				Ratelimit    float64 `json:"ratelimit"`
				Toolarge     float64 `json:"toolarge"`
				Unauthorized float64 `json:"unauthorized"`
				Unavailable  float64 `json:"unavailable"`
				Validate     float64 `json:"validate"`
			} `json:"errors"`
			Valid struct {
				Accepted    float64 `json:"accepted"`
				Count       float64 `json:"count"`
				Notmodified float64 `json:"notmodified"`
				Ok          float64 `json:"ok"`
			} `json:"valid"`
		} `json:"response"`
		Unset float64 `json:"unset"`
	} `json:"sourcemap"`
}
