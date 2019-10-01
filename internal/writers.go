package papersave

import (
	"path"
	"fmt"
	"os"
)

func getWriterFH(path string) (fh *os.File, err error) {
	fh, err = os.Create(path)
	Panicp(err)
	return
}

func (self PaperSave) WriteText() {
	err := self.writeTemplate(
		"templates/papersave.txt",
		os.Stdout)
	Panicp(err)

}


func copyBinFile(src, dir string) error {
	content, err := getFile(src)
	if err != nil {
		return err
	}
	out, err := os.Create(fmt.Sprintf("%s/%s", dir, path.Base(src)))
	defer out.Close()
	if err != nil {
		return err
	}
	_, err = out.Write(content)
	if err != nil {
		return err
	}
	return nil
}

func copyTexFiles(workDir string) {
	files := []string {
		"iosevka-light.ttf",
		"iosevka-lightitalic.ttf",
		"iosevka-semibold.ttf",
		"iosevka-semibolditalic.ttf"}

	for f := range files {
		err := copyBinFile(fmt.Sprintf("/templates/%s", files[f]), workDir)
		if err != nil {
			panic(err)
		}
	}
}	


func (self PaperSave) WritePDF() {
	self.GenQRCode()
	outPath :=  fmt.Sprintf("%s/%s.tex", self.WorkDir, self.Filename)
	fh, err := os.Create(outPath)
	Panicp(err)
	defer fh.Close()
	err = self.writeTemplate(
		"templates/papersave.tex",
		fh)
	Panicp(err)
	copyTexFiles(self.WorkDir)
	for i:=0; i<2; i++ {
		runCommand(self.WorkDir,
			"xelatex", "-halt-on-error",
			//"-interaction=batchmode",
		fmt.Sprintf("%s.tex", self.Filename))
	}
	os.Rename(fmt.Sprintf("%s/%s.pdf", self.WorkDir, self.Filename),
		fmt.Sprintf("%s.pdf", self.Filename))
	if !self.Options.Keep {
		os.RemoveAll(self.WorkDir)
	}
}


func (self PaperSave) writeTemplate(tplFile string, out *os.File) error {
	template, err := getTemplate(tplFile)
	Panicp(err)
	err = template.Execute(out, self)
	Panicp(err)
	return err
}

