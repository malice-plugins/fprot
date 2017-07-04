package main

const tpl = `#### F-PROT
{{- with .Results }}
| Infected      | Result      | Engine      | Updated      |
|:-------------:|:-----------:|:-----------:|:------------:|
| {{.Infected}} | {{.Result}} | {{.Engine}} | {{.Updated}} |
{{ end -}}
`

// func printMarkDownTable(fprot FPROT) {

// 	fmt.Println("#### F-PROT")
// 	table := clitable.New([]string{"Infected", "Result", "Engine", "Updated"})
// 	table.AddRow(map[string]interface{}{
// 		"Infected": fprot.Results.Infected,
// 		"Result":   fprot.Results.Result,
// 		"Engine":   fprot.Results.Engine,
// 		"Updated":  fprot.Results.Updated,
// 	})
// 	table.Markdown = true
// 	table.Print()
// }
