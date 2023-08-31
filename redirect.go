package oakmux

import (
	"net/http"
	"strings"
)

type redirect struct {
	location   string
	statusCode int
}

func (h *redirect) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	http.Redirect(w, r, h.location, h.statusCode)
	return nil
}

func NewTemporaryRedirect(to string) Handler {
	return &redirect{
		location:   to,
		statusCode: http.StatusTemporaryRedirect,
	}
}

func NewPermanentRedirect(to string) Handler {
	return &redirect{
		location:   to,
		statusCode: http.StatusPermanentRedirect,
	}
}

func (o *options) injectTrailingSlashRedirects() (err error) {
	if !o.redirectToTrailingSlash && !o.redirectFromTrailingSlash {
		return nil // nothing to redirect
	}

	return o.tree.Walk(func(n *Node) (ok bool, err error) {
		if n.Leaf != nil && n.TrailingSlashLeaf == nil && o.redirectToTrailingSlash {
			path := n.Leaf.String()
			// fmt.Println("###", path, strings.TrimPrefix(path, "/"+o.prefix))
			// fmt.Println(n.Leaf.name + ":slashRedirect")
			// fmt.Println(NewTemporaryRedirect(path))
			if err = WithRouteHandler(
				n.Leaf.name+":slashRedirect",
				strings.TrimPrefix(path, "/"+o.prefix)+"/",
				NewTemporaryRedirect(path),
			)(o); err != nil {
				return false, err
			}
		}
		// fmt.Println("===============================")

		if n.TrailingSlashLeaf != nil && n.Leaf == nil && o.redirectFromTrailingSlash {
			path := n.TrailingSlashLeaf.String()
			if err = WithRouteHandler(
				n.TrailingSlashLeaf.name+":slashRedirect",
				path[len(o.prefix)+1:len(path)-1], // strip
				NewTemporaryRedirect(path),
			)(o); err != nil {
				return false, err
			}
		}

		return true, nil
	})
}
