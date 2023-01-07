package rest

import (
	"fmt"
	"net/http"

	"github.com/delveper/mylib/app/ent"
	"github.com/go-chi/chi/v5"
)

type Reader struct{ ReaderLogic }

func NewReader(logic ReaderLogic, logger ent.Logger) Reader {
	return Reader{ReaderLogic: logic}
}

func (r Reader) Route(router chi.Router) {
	router.Use()
	router.Method(http.MethodPost, "/readers", r.Create())
}

func (r Reader) Create() HandlerLoggerFunc {
	return func(rw http.ResponseWriter, req *http.Request, logger ent.Logger) {
		logger.Println("o la la")
		fmt.Fprintln(rw, "bla bla")
		/*		var reader ent.Reader
				if err := r.decodeBody(&reader); err != nil {
					r.Write(http.StatusBadRequest, Message{Message: MsgBadRequest, Details: err.Error()})
					r.Errorf("Failed decoding reader data from request.", "request", req, "error", err)

					return
				}

				if err := reader.Validate(); err != nil {
					r.Write(http.StatusBadRequest, Message{Message: MsgBadRequest, Details: err.Error()})
					r.Debugf("Failed validating reader: %v", err)

					return
				}

				err := r.SignUp(context.Background(), reader)
				if err != nil {
					switch {
					case errors.Is(err, exc.ErrDuplicateEmail):
						r.Write(http.StatusConflict, Message{Message: MsgConflict, Details: exc.ErrDuplicateEmail.Error()})
					case errors.Is(err, exc.ErrDuplicateID):
						r.Write(http.StatusConflict, Message{Message: MsgConflict, Details: exc.ErrDuplicateID.Error()})
					default:
						r.Write(http.StatusInternalServerError, Message{Message: MsgInternalSeverErr})
					}

					r.Errorf("Failed creating reader: %+v", err)

					return
				}

				r.Write(http.StatusCreated, Message{Message: MsgSuccess})
				r.Debugw("Reader successfully created")

		*/
	}
}
