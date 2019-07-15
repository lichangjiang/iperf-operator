package util

import (
	"testing"

	"gotest.tools/assert"
)

func TestSendEmail(t *testing.T) {

	str := `
    <table>
        <tr>
            <th>col1</th>
            <th>col2</th>
        </tr>
        <tr>
            <th>hello world</th>
            <th>test<br>test<br></th>
        </tr>
    </table>
    `
	err := SendEmail("305120108@qq.com", "test", str)
	assert.NilError(t, err)
}
