# Functionality

There's three roles to access the site and 4 backend roles

* Unverified (1) Has not confirmed email
* User (2) Two possible states, profile verified, profile not verified
* Certifier (3) Can create courses and create certificates to send to any account
* Admin (4) CRM of entire site

## User signup

A "user" can exist in two states, they have recieved a certificate without logging in or signing up to the application. 
This is to be able to send certificates to non users. The flow for a user that was sent a certificate is the same as a new user sign up, the
only difference is that the user sent a certificate will recive an email to create an account.

* An email is sent on sucessful form submit / User certificate sent
* User must confirm email with uid (click link) to be able to signin

On first signin the user will need to submit a account profile verification (kyc), they will also attach 3 images, this will be sent to the Admin for.
manual verification. The user will not be able to view pages unless their account has been verified.

Once verified they will have access to the pages.

## User signin

* User is routed to their specific page under their account role.

All users except unverified have access to signing out and profile info.

## Role (2): User

* Ability to view certificates and download certificates.

## Role (3): Certifier

* Create course
* View Course information / anylitics (How many certificates were sent out, how many were updated, etc)
* Create certificate (Sends certificate to email)
* Course and Certificate CRUD

## Role (4): Admin

* User Crud
* Assign User roles
* User profile verification.
* View Course / certificates for users Anylitics (For billing purposes)