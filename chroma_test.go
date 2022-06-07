package boondoggle

import (
	"testing"
)

var chromaInput = "```sql\n" +
	"SELECT * FROM TABLE;\n" +
	"```\n"

var chromaOutput = `<pre tabindex="0" class="chroma"><code><span class="line"><span class="cl"><span class="k">SELECT</span><span class="w"> </span><span class="o">*</span><span class="w"> </span><span class="k">FROM</span><span class="w"> </span><span class="k">TABLE</span><span class="p">;</span><span class="w">
</span></span></span></code></pre>\n`

func TestChroma(t *testing.T) {
	var example Article
	example.Raw = []byte(chromaInput)

	if err := ChromaCode(&example); err != nil {
		t.Fatal(err)
	}

	if string(example.Raw) != chromaOutput {
		t.Errorf("unexpected chroma output: %s", example.Raw)
	}
}
