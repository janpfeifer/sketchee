// main(): called in body.onLoad()
function main() {
    // console.log("main()");
    StartWASM();
}

function StartWASM() {
    // console.log("StartWASM()");
    if (!WebAssembly.instantiateStreaming) {
        // console.log("StartWASM(): create instantiateStreaming");
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }

    const go = new Go();
    let mod, inst;

    async function run() {
        // console.log("StartWASM().run()");
        await go.run(inst)
        // console.log("StartWASM().run() finished");
    }

    WebAssembly.instantiateStreaming(
        fetch("main.wasm", {"cache": "no-cache"}),
        go.importObject).then(
        (result) => {
            // console.log("StartWASM(): fetched main.wasm");
            inst = result.instance;
            mod = result.module;
            run().then((_) => {
                console.log("main.wasm: returned from run!?")
            });
        }).catch((err) => {
        console.error("Failed to fetch main.wasm: " + err);
    });
}
