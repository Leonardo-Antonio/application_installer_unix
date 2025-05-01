# Application installer

Instala aplicaciones de forma sencilla y rápida.


## Instalación
```bash
sudo mkdir -p /usr/local/bin/application_installer && find /usr/local/bin/application_installer -type f -name "linux*" | xargs sudo rm && sudo wget -P /usr/local/bin/application_installer https://github.com/Leonardo-Antonio/application_installer_unix/releases/download/v0.0.2-dev/linux && sudo chmod 0777 /usr/local/bin/application_installer/linux && echo 'alias app_installer="sudo /usr/local/bin/application_installer/linux"' >> ~/.zshrc && source ~/.zshrc && echo "Application installer installed successfully (execute in terminal: app_installer)" && app_installer
```

## Uso
```bash
app_installer
```