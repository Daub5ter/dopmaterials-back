package handlers

import (
	"api-gateway/internal/models/consts"
	"api-gateway/internal/models/errs"
	"api-gateway/internal/tools/server/servparams"
	"api-gateway/pkg/code"
	"io"
	"net/http"
)

func GetMaterials() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		categoryId := params.Get(consts.CategoryIdParam)
		findPart := params.Get(consts.FindPartParam)
		offset := params.Get(consts.OffsetParam)
		request, err := http.NewRequest(
			"GET",
			consts.ContentUrl+"materials?"+
				consts.CategoryIdParam+"="+categoryId+
				"&"+consts.FindPartParam+"="+findPart+
				"&"+consts.OffsetParam+"="+offset,
			nil,
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = response.Body.Close() }()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
		_, _ = w.Write(body)
	}
}

func GetMaterial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := servparams.ViewParam(r, consts.IdParam)
		if id == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		request, err := http.NewRequest(
			"GET",
			consts.ContentUrl+"materials/"+id,
			nil,
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = response.Body.Close() }()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
		_, _ = w.Write(body)
	}
}

func AddMaterial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := http.NewRequest(
			"POST",
			consts.ContentUrl+"materials",
			r.Body,
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = response.Body.Close() }()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
		_, _ = w.Write(body)
	}
}

func UpdateMaterial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := http.NewRequest(
			"PUT",
			consts.ContentUrl+"materials",
			r.Body,
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = response.Body.Close() }()

		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
	}
}

func DeleteMaterial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := servparams.ViewParam(r, consts.IdParam)
		if id == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		request, err := http.NewRequest(
			"DELETE",
			consts.ContentUrl+"materials/"+id,
			nil,
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = response.Body.Close() }()

		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
	}
}

func RestoreMaterial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := servparams.ViewParam(r, consts.IdParam)
		if id == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		request, err := http.NewRequest(
			"PATCH",
			consts.ContentUrl+"materials/"+id+"/restore",
			nil,
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = response.Body.Close() }()

		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
	}
}

func GetMainCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := http.NewRequest(
			"GET",
			consts.ContentUrl+"categories",
			nil,
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = response.Body.Close() }()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
		_, _ = w.Write(body)
	}
}

func GetSubsidiariesCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := servparams.ViewParam(r, consts.IdParam)
		if id == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		request, err := http.NewRequest(
			"GET",
			consts.ContentUrl+"categories/"+id+"/subsidiaries",
			nil,
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = response.Body.Close() }()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
		_, _ = w.Write(body)
	}
}

func GetCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := servparams.ViewParam(r, consts.IdParam)
		if id == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		request, err := http.NewRequest(
			"GET",
			consts.ContentUrl+"categories/"+id,
			nil,
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = response.Body.Close() }()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
		_, _ = w.Write(body)
	}
}

func AddCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := http.NewRequest(
			"POST",
			consts.ContentUrl+"categories",
			r.Body,
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = response.Body.Close() }()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
		_, _ = w.Write(body)
	}
}

func UpdateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := http.NewRequest(
			"PUT",
			consts.ContentUrl+"categories",
			r.Body,
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = response.Body.Close() }()

		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
	}
}

func DeleteCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := servparams.ViewParam(r, consts.IdParam)
		if id == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		request, err := http.NewRequest(
			"DELETE",
			consts.ContentUrl+"categories/"+id,
			nil,
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = response.Body.Close() }()

		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
	}
}
