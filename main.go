package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

type format string

const (
	APPIMAGE format = "appimage"
	DEB      format = "deb"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	var rootCmd = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			var format string
			promptformat := &survey.Select{
				Message: "Formato de la aplicacion: ",
				Options: []string{string(DEB), string(APPIMAGE)},
			}
			err := survey.AskOne(promptformat, &format)
			if err != nil {
				fmt.Println("Error al seleccionar el formato:", err)
				return
			}

			var path string
			promptpath := &survey.Input{
				Message: fmt.Sprintf("Donde esta la aplicacion (%s)? [~/Descargas/discord.tar.gz]", format),
			}
			err = survey.AskOne(promptpath, &path)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			if path == "" {
				fmt.Println("No se especifico la ruta de la aplicacion")
				os.Exit(1)
			}

			if !strings.Contains(strings.ToLower(path), strings.ToLower(format)) {
				fmt.Printf("La ruta de la aplicacion no contiene el formato seleccionado (%s) \n", format)
				os.Exit(1)
			}

			if strings.Contains(strings.ToLower(path), "~/") {
				path = strings.Replace(path, "~/", usr.HomeDir+"/", 1)
			}
			file, err := os.OpenFile(path, os.O_RDONLY, 0644)
			if os.IsNotExist(err) {
				fmt.Println("La aplicacion no existe en la ruta especificada: " + path)
				os.Exit(1)
			}

			if err != nil {
				fmt.Println("Error al abrir la aplicacion:", err)
				os.Exit(1)
			}
			defer file.Close()

			var nameApp string
			promptNameApp := &survey.Input{
				Message: "Nombre que le quieres dar a la aplicación? [Discord]",
			}
			err = survey.AskOne(promptNameApp, &nameApp)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			filepathapplication := filepath.Join(usr.HomeDir, ".config_application_installer", strings.ToLower(nameApp))
			if err := os.MkdirAll(filepathapplication, 0755); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			fmt.Println("Configuración creada en: " + filepathapplication)

			executablepath := filepath.Join(filepathapplication, nameApp+"."+format)
			if err := os.Rename(path, executablepath); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			install(format, executablepath)
		},
	}

	rootCmd.Execute()
}

func install(format string, path string) {
	switch format {
	case string(DEB):
		fmt.Println("Instalando...")
		cmd := exec.Command("sudo", "dpkg", "-i", path)

		var out bytes.Buffer
		var stderr bytes.Buffer

		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error al ejecutar el comando: %v \n", err)
		}

		fmt.Println("Instalación finalizada")

	case string(APPIMAGE):
		cmd := exec.Command("chmod", "a+x", path)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error al dar permisos de ejecución: %v \n", err)
		}

		// Ejecutar el archivo AppImage
		cmd = exec.Command(path)
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error al ejecutar el archivo AppImage: %v \n", err)
		}
	default:
		fmt.Println("Formato no soportado")
	}
}
