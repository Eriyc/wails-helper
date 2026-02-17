package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Eriyc/wails-helper"
	"github.com/Eriyc/wails-helper/pkg/providers/autoupdate/github"
	licensecustom "github.com/Eriyc/wails-helper/pkg/providers/license/custom"
)

func main() {
	// Example: Setting up the Wails helper with GitHub auto-updates
	// and a custom license server

	// Configure auto-update provider for GitHub releases
	updateProvider := github.NewProvider(github.Config{
		Owner: "your-username",
		Repo:  "your-app-repo",
		Token: "", // Optional: use token for private repos or higher rate limits
	})

	// Configure license provider for custom server
	licenseProvider := licensecustom.NewProvider(licensecustom.Config{
		BaseURL: "https://license.yourdomain.com",
		APIKey:  "your-api-key",
	})

	// Create the helper instance
	h := helper.New(helper.Config{
		AppName:         "My Awesome App",
		AppVersion:      "1.0.0",
		UpdateProvider:  updateProvider,
		LicenseProvider: licenseProvider,
		// Note: In a real Wails3 app, you would also set:
		// WindowManager:  wails3WindowManager,
		// DialogManager:  wails3DialogManager,
		// EventManager:   wails3EventManager,
	})

	ctx := context.Background()

	// Example 1: Check for updates
	fmt.Println("Checking for updates...")
	updateInfo, err := h.CheckForUpdates(ctx)
	if err != nil {
		log.Printf("Error checking for updates: %v", err)
	} else {
		fmt.Printf("Current version: %s\n", updateInfo.CurrentVersion)
		if updateInfo.Available {
			fmt.Printf("New version available: %s\n", updateInfo.LatestVersion.Version)
			fmt.Printf("Release notes: %s\n", updateInfo.LatestVersion.Notes)
			fmt.Printf("Download URL: %s\n", updateInfo.LatestVersion.DownloadURL)

			// Optionally download the update
			// reader, err := h.Updates().DownloadUpdate(ctx, updateInfo.LatestVersion)
			// if err != nil {
			//     log.Fatal(err)
			// }
			// defer reader.Close()
			// ... handle the downloaded update
		} else {
			fmt.Println("You're up to date!")
		}
	}

	// Example 2: Validate a license
	fmt.Println("\nValidating license...")
	licenseKey := "EXAMPLE-LICENSE-KEY-1234"
	result, err := h.ValidateLicense(ctx, licenseKey)
	if err != nil {
		log.Printf("Error validating license: %v", err)
	} else {
		if result.Valid {
			fmt.Println("License is valid!")
			if result.License != nil {
				fmt.Printf("License type: %s\n", result.License.Type)
				fmt.Printf("Features: %v\n", result.License.Features)
				if result.License.ExpiresAt != nil {
					fmt.Printf("Expires at: %s\n", result.License.ExpiresAt)
				}
			}
		} else {
			fmt.Printf("License is invalid: %s\n", result.Error)
		}
	}

	// Example 3: Check if a specific feature is enabled
	fmt.Println("\nChecking for premium feature...")
	hasPremium, err := h.License().CheckFeature(ctx, licenseKey, "premium")
	if err != nil {
		log.Printf("Error checking feature: %v", err)
	} else {
		if hasPremium {
			fmt.Println("Premium features are enabled!")
		} else {
			fmt.Println("Premium features are not available with this license.")
		}
	}

	// Example 4: Get version history
	fmt.Println("\nFetching version history...")
	versions, err := h.Updates().GetVersionHistory(ctx, 5)
	if err != nil {
		log.Printf("Error fetching version history: %v", err)
	} else {
		fmt.Printf("Found %d recent versions:\n", len(versions))
		for i, version := range versions {
			fmt.Printf("%d. v%s - %s\n", i+1, version.Version, version.ReleaseDate)
		}
	}

	// In a real Wails3 application, you would also use:
	//
	// // Show dialogs
	// h.ShowInfo(ctx, "Welcome", "Welcome to My Awesome App!")
	//
	// // Create windows
	// win, _ := h.Windows().Create(ctx, window.Options{
	//     Title:  "Main Window",
	//     Width:  1024,
	//     Height: 768,
	// })
	//
	// // Register keyboard shortcuts
	// h.Events().RegisterShortcut(ctx, events.Shortcut{
	//     Key:       "s",
	//     Modifiers: events.ModifierControl,
	//     Handler: func(event interface{}) error {
	//         fmt.Println("Ctrl+S pressed!")
	//         return nil
	//     },
	// })
}
