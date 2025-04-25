package loader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/henrypekny/runt/fonts" // for embedded VT323
)

var (
	loaderPaths []string
	fontCache   = make(map[string]rl.Font)
	texCache    = make(map[string]rl.Texture2D)
	mu          sync.Mutex
)

func init() {
	// 1) location of the running binary
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		loaderPaths = append(loaderPaths,
			dir,
			filepath.Join(dir, "assets"),
		)
		if runtime.GOOS == "darwin" {
			// inside a .app bundle
			loaderPaths = append(loaderPaths,
				filepath.Join(dir, "..", "Resources"),
			)
		}
	}
	// 2) development fallbacks
	loaderPaths = append(loaderPaths,
		".",
		"assets",
		"../assets",
		"../../runt/assets",
		"../runt/fonts",
		"../../runt/fonts",
	)
}

// Resolve searches each loaderPaths entry for `path`.
func Resolve(path string) (string, error) {
	for _, dir := range loaderPaths {
		full := filepath.Join(dir, path)
		if fi, err := os.Stat(full); err == nil && !fi.IsDir() {
			return full, nil
		}
	}
	return "", fmt.Errorf("runt: asset %q not found", path)
}

// LoadFont loads (and caches) a font at the given size, using disk or embedded VT323.
func LoadFont(path string, size int32) rl.Font {
	key := fmt.Sprintf("%s#%d", path, size)

	mu.Lock()
	defer mu.Unlock()
	if f, ok := fontCache[key]; ok {
		return f
	}

	// 1) try on-disk
	full, err := Resolve(path)

	// 2) embedded fallback for VT323
	if err != nil && filepath.Base(path) == "VT323-Regular.ttf" {
		tmp, e2 := ioutil.TempFile("", "vt323-*.ttf")
		if e2 != nil {
			panic(fmt.Errorf("loader: cannot create temp for embedded VT323: %w", e2))
		}
		defer tmp.Close()
		if _, e3 := tmp.Write(fonts.VT323TTF); e3 != nil {
			panic(fmt.Errorf("loader: cannot write embedded VT323: %w", e3))
		}
		full = tmp.Name()
	}
	if full == "" {
		panic(err)
	}

	// 3) load + point-filter
	fnt := rl.LoadFontEx(full, size, nil, 0)
	rl.SetTextureFilter(fnt.Texture, rl.FilterPoint)
	fontCache[key] = fnt
	return fnt
}

// LoadTexture loads (and caches) a Texture2D, forces point-filtering.
func LoadTexture(path string) rl.Texture2D {
	mu.Lock()
	defer mu.Unlock()
	if t, ok := texCache[path]; ok {
		return t
	}

	// 1) try on-disk resolve
	full, err := Resolve(path)
	if err != nil {
		// also try basename only
		full, err = Resolve(filepath.Base(path))
	}
	if err != nil {
		panic(err)
	}

	// 2) load + point-filter
	tex := rl.LoadTexture(full)
	rl.SetTextureFilter(tex, rl.FilterPoint)
	texCache[path] = tex
	return tex
}
