# Cloudflare Tunnel (Zero Trust) Guide

This guide explains how to expose a local application to the public internet using Cloudflare Tunnel, part of Cloudflare's Zero Trust platform.

## Prerequisites

- You need a domain name that you own and can manage its DNS.
- You need a Cloudflare account with your domain added to it.

(See the "Alternatives" section below if you do not have a domain).

## Step 1: Install `cloudflared`

`cloudflared` is the command-line tool for Cloudflare Tunnel.

### For Windows

1.  **Download `cloudflared`** from the [Cloudflare downloads page](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/downloads/).
2.  Create a new folder, for example `C:\Cloudflared`.
3.  Move the downloaded file into this new directory and rename it to `cloudflared.exe`.
4.  Open a command prompt (CMD) or PowerShell and navigate to the directory where you saved `cloudflared.exe`.
5.  Verify the installation by running:
    ```powershell
    .\cloudflared.exe --version
    ```

### For Linux (Debian/Ubuntu)

1.  **Download the package:**
    ```bash
    wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
    ```

2.  **Install the package:**
    ```bash
    sudo dpkg -i cloudflared-linux-amd64.deb
    ```
    If you encounter dependency issues, run:
    ```bash
    sudo apt-get install -f
    ```

## Step 2: Authenticate `cloudflared`

This command will open a browser window to log in to your Cloudflare account and authorize the tool.

```bash
cloudflared tunnel login
```

## Step 3: Create a Tunnel

Create a tunnel and give it a memorable name.

```bash
cloudflared tunnel create <tunnel-name>
```
For example:
```bash
cloudflared tunnel create ideb-app
```
This will create a credentials file (e.g., `~/.cloudflared/TUNNEL-ID.json`).

## Step 4: Configure the Tunnel

Create a configuration file to tell `cloudflared` where to route traffic.

1.  Create a `config.yml` file in your `~/.cloudflared` directory.
2.  Add the following content, adjusting the placeholders:

    ```yaml
    tunnel: <tunnel-name>
    credentials-file: /path/to/your/.cloudflared/<YOUR_TUNNEL_ID>.json

    ingress:
      - hostname: <your-subdomain>.<your-domain>.com
        service: http://localhost:8080
      - service: http_status:404
    ```

    **Important:**
    - Replace `<tunnel-name>` with the name you chose in the previous step.
    - Replace `/path/to/your/.cloudflared/<YOUR_TUNNEL_ID>.json` with the actual path to the credentials file.
    - Replace `<your-subdomain>.<your-domain>.com` with the public URL you want to use.

## Step 5: Route DNS to the Tunnel

Link your chosen hostname to the tunnel.

```bash
cloudflared tunnel route dns <tunnel-name> <your-subdomain>.<your-domain>.com
```
For example:
```bash
cloudflared tunnel route dns ideb-app ideb.yourdomain.com
```

## Step 6: Run the Tunnel

Start the tunnel. Your application will be accessible at the hostname you configured.

```bash
cloudflared tunnel run <tunnel-name>
```

## Zero Trust Security

With the tunnel running, you can enhance security by adding access policies.

1.  Go to the **Zero Trust** section of your Cloudflare dashboard.
2.  Navigate to **Access -> Applications**.
3.  Add a new application and select "Self-hosted".
4.  Configure the application, pointing it to the hostname you created.
5.  Create policies to define who can access your application (e.g., require a specific email domain or a one-time PIN).

---

## Alternatives for Development (If You Don't Have a Domain)

If you don't have your own domain, you can use one of these options for development and testing.

### Option 1: Cloudflare Quick Tunnels

For quick testing, Cloudflare offers a feature that works just like ngrok. It gives you a random, temporary hostname.

Simply run this single command:

```bash
cloudflared tunnel --url http://localhost:8080
```

This will connect and print a public URL ending in `.trycloudflare.com` that you can use immediately.

### Option 2: Get a Free Domain

You can get a free domain from a service like **Freenom** ([https://www.freenom.com](https://www.freenom.com)). They offer domains with extensions like `.tk`, `.ml`, `.ga`, `.cf`, and `.gq`.

**Important Note:** While these domains work for development and testing, they can be unreliable and are not recommended for important production applications.

Once you have a free domain, you can add it to your Cloudflare account and follow the main guide above.
