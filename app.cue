package main

resources: {
	(#StaticWebApp & {
		_name: "acald"
	}).resources
}

#StaticWebApp: {
	_name: string
	_documents: {
		root:     string | *"website"
		index:    string | *"index.html"
		error404: string | *"404.html"
		json:     string | *"resume.json"
		js:       string | *"main.js"
	}

	resources: {
		"\(_name)-group": {
			type: "azure-native:resources:ResourceGroup"
			properties: {
				location: "EastUS"
			}
		}

		"\(_name)storage": {
			type: "azure-native:storage:StorageAccount"
			properties: {
				resourceGroupName: "${\(_name)-group.name}"
				kind:              "StorageV2"
				sku: {
					name: "Standard_LRS"
				}
			}
		}

		"\(_name)-website": {
			type: "azure-native:storage:StorageAccountStaticWebsite"
			properties: {
				resourceGroupName: "${\(_name)-group.name}"
				accountName:       "${\(_name)storage.name}"
				indexDocument:     "\(_documents.index)"
				error404Document:  "\(_documents.error404)"
			}
		}

		"index.html": {
			type: "azure-native:storage:Blob"
			properties: {
				resourceGroupName: "${\(_name)-group.name}"
				accountName:       "${\(_name)storage.name}"
				containerName:     "${\(_name)-website.containerName}"
				contentType:       "text/html"
				type:              "Block"
				source: {
					"Fn::FileAsset": "./\(_documents.root)/\(_documents.index)"
				}
			}
		}

		"404.html": {
			type: "azure-native:storage:Blob"
			properties: {
				resourceGroupName: "${\(_name)-group.name}"
				accountName:       "${\(_name)storage.name}"
				containerName:     "${\(_name)-website.containerName}"
				contentType:       "text/html"
				type:              "Block"
				source: {
					"Fn::FileAsset": "./\(_documents.root)/\(_documents.index)"
				}
			}
		}

		"favicon.png": {
			type: "azure-native:storage:Blob"
			properties: {
				resourceGroupName: "${\(_name)-group.name}"
				accountName:       "${\(_name)storage.name}"
				containerName:     "${\(_name)-website.containerName}"
				contentType:       "image/png"
				type:              "Block"
				source: {
					"Fn::FileAsset": "./\(_documents.root)/\(_documents.index)"
				}
			}
		}

		"resume.json": {
			type: "azure-native:storage:Blob"
			properties: {
				resourceGroupName: "${\(_name)-group.name}"
				accountName:       "${\(_name)storage.name}"
				containerName:     "${\(_name)-website.containerName}"
				contentType:       "application/json"
				type:              "Block"
				source: {
					"Fn::FileAsset": "./\(_documents.root)/\(_documents.json)"
				}
			}
		}

		"main.js": {
			type: "azure-native:storage:Blob"
			properties: {
				resourceGroupName: "${\(_name)-group.name}"
				accountName:       "${\(_name)storage.name}"
				containerName:     "${\(_name)-website.containerName}"
				contentType:       "text/javascript"
				type:              "Block"
				source: {
					"Fn::FileAsset": "./\(_documents.root)/\(_documents.js)"
				}
			}
		}

		"cdnprofile": {
			type: "azure-native:cdn:Profile"
			properties: {
				location:          "EastUS"
				profileName:       "resumeprofile"
				resourceGroupName: "${\(_name)-group.name}"
				sku: {
					name: "Standard_Verizon"
				}
			}
		}

		"endpointcreation": {
			type: "azure-native:cdn:Endpoint"
			properties: {
				contentTypesToCompress: [
					"text/plain",
					"text/html",
					"text/css",
					"text/javascript",
					"application/x-javascript",
					"application/javascript",
					"application/json",
					"application/xml",
				]
				endpointName:         "personal-website"
				isCompressionEnabled: true
				isHttpAllowed:        true
				isHttpsAllowed:       true
				location:             "Global"
				originHostHeader:     "acaldstorage59084b7e.z13.web.core.windows.net"
				origins: [{
					enabled:   true
					hostName:  "acaldstorage59084b7e.z13.web.core.windows.net"
					httpPort:  80
					httpsPort: 443
					name:      "acaldstorage59084b7e-z13-web-core-windows-net"
				}]
				profileName:                "resumeprofile"
				queryStringCachingBehavior: "IgnoreQueryString"
				resourceGroupName:          "${\(_name)-group.name}"
			}
			options: {
				protect: true
			}
		}

		"resume-acaldwell-dev": {
			type: "azure-native:cdn:CustomDomain"
			properties: {
				customDomainName:  "resume-acaldwell-dev"
				endpointName:      "personal-website"
				hostName:          "resume.acaldwell.dev"
				profileName:       "resumeprofile"
				resourceGroupName: "${\(_name)-group.name}"
			}
			options: {
				protect: true
			}
		}
	}
}
