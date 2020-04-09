package nylas

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

// File represents a file in the Nylas system.
type File struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	AccountID string `json:"account_id"`

	ContentType string `json:"content_type"`
	Filename    string `json:"filename"`
	Size        int    `json:"size"`
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// UploadFile uploads a file to be used as an attachment.
// See: https://docs.nylas.com/reference#upload
func (c *Client) UploadFile(
	ctx context.Context, filename string, file io.Reader,
) (File, error) {
	g, ctx := errgroup.WithContext(ctx)
	req, err := c.newUserRequest(ctx, http.MethodPost, "/files", nil)
	if err != nil {
		return File{}, err
	}

	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	req.Body = r
	req.Header.Set("Content-Type", m.FormDataContentType())

	g.Go(func() (err error) {
		defer w.Close() // nolint: errcheck
		defer m.Close() // nolint: errcheck

		contentType := mime.TypeByExtension(filepath.Ext(filename))
		if contentType == "" {
			contentType = "application/octet-stream" // fallback
		}

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="file"; filename="%s"`,
				escapeQuotes(filename)))
		h.Set("Content-Type", contentType)
		part, err := m.CreatePart(h)
		if err != nil {
			return err
		}

		_, err = io.Copy(part, file)
		return err
	})

	var resp []File
	g.Go(func() error {
		if err := c.do(req, &resp); err != nil {
			return err
		} else if len(resp) == 0 {
			return errors.New("no file returned")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return File{}, err
	}
	return resp[0], nil
}
