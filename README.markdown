Shync - Shopify Email Template Syncer
=====================================

Shync can download the Shopify email templates from your Shopify store admin to your local machine,
and upload (or sync) your email templates on your local machine to Shopify.

How to use
==========

```
# Download templates from Shopify to . directory
% shync download .

# Upload templates from . to Shopify
% shync sync .
```

Configuration
=============

Shync needs to know some information to do it's job.
You can set them using following environment variables, or in a .env file.

SHYNC_STORE: The Shopify store URL. eg: `https://{shopname}.myshopify.com`. Required.

SHYNC_USERNAME: The Shopify admin username. Required.

SHYNC_PASSWORD: The Shopify admin password. Required.
