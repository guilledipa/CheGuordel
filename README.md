# Che Guordel

"Che Guordel!" es una implementacion del juego Wordle, en Go, para un publico
Argentino.

Sera desarrollado usando https://ebitengine.org

## Build

### Linux

```shell
env GOOS=js GOARCH=wasm go build -o yourgame.wasm github.com/guilledipa/cheguordel
```

```shell
cp $(go env GOROOT)/misc/wasm/wasm_exec.js .
```

### Windows

```shell
$Env:GOOS = 'js'
$Env:GOARCH = 'wasm'
go build -o yourgame.wasm github.com/guilledipa/cheguordel
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
```

```shell
$goroot = go env GOROOT
cp $goroot\misc\wasm\wasm_exec.js .
```

## HTML

Create this HTML:

```html
<!DOCTYPE html>
<script src="wasm_exec.js"></script>
<script>
// Polyfill
if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

const go = new Go();
WebAssembly.instantiateStreaming(fetch("yourgame.wasm"), go.importObject).then(result => {
    go.run(result.instance);
});
</script>
```

Luego abrir el HTML en tu navegador.

Para embeber el juego en una pagina web se recomienda usar un iframe.

Si el HTML visto mas arriba se llama `main.html` el HTML host deberia ser:

```html
<!DOCTYPE html>
<iframe src="main.html" width="640" height="480"></iframe>
```

Si ves este mensaje en Chrome:

```none
The AudioContext was not allowed to start. It must be resume (or created) after a user gesture on the page. https://goo.gl/7K7WLu
```

Se resuelve agregando `allow="autoplay"` en el iframe.
