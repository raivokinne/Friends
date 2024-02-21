// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package register

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/RaivoKinne/Chatify/web/templates/layout"
)

func Register(err string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex items-center justify-center h-screen\"><form class=\"grid grid-cols-1 gap-6 w-[500px] px-8 pt-6 pb-8 mb-4 bg-white rounded shadow-md place-items-center\" method=\"POST\" action=\"/register\"><div class=\"grid w-[400px]\"><label for=\"username\" class=\"block mb-2 text-sm font-bold text-gray-700\">Username</label> <input id=\"username\" name=\"username\" type=\"text\" placeholder=\"Username\" class=\"w-full px-3 py-2 leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline\"></div><div class=\"grid w-[400px]\"><label for=\"email\" class=\"block mb-2 text-sm font-bold text-gray-700\">Email</label> <input id=\"email\" name=\"email\" type=\"text\" placeholder=\"Email\" class=\"w-full px-3 py-2 leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline\"></div><div class=\"grid w-[400px]\"><label for=\"password\" class=\"block mb-2 text-sm font-bold text-gray-700\">Password</label> <input id=\"password\" name=\"password\" type=\"password\" placeholder=\"Password\" class=\"w-full px-3 py-2 mb-3 leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline\"></div><div class=\"grid w-[400px]\"><label for=\"confirmPassword\" class=\"block mb-2 text-sm font-bold text-gray-700\">Confirm Password</label> <input id=\"confirmPassword\" name=\"confirmPassword\" type=\"password\" placeholder=\"Confirm Password\" class=\"w-full px-3 py-2 mb-3 leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline\"></div><div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if err != "" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p class=\"text-red-500\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(err)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/auth/register/register.templ`, Line: 36, Col: 35}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div><button type=\"submit\" class=\"px-4 py-2 font-bold text-white bg-blue-500 rounded hover:bg-blue-700 focus:outline-none focus:shadow-outline w-[400px]\">Register</button></div></form></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = layout.BaseLayout("Register").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}