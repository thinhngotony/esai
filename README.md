# Project: Bill Splitting Automation Web Application

## Objective

Develop a full-stack web application to automate the bill-splitting process. Users should be able to:

* Upload or input itemized bills
* Assign items to different participants
* Apply discounts and other fees
* Calculate each participant's share, rounding up as needed

## Features

### Frontend

* **Bill Input Form**: Provide an interface where users can manually enter items with names and prices or upload an image of a bill. For image uploads, use OCR (Optical Character Recognition) to extract item names and prices.
* **Participant Assignment**: Allow users to assign each item to a specific participant. Each participant can be tagged with their item selections.
* **Discounts and Fees**: Include fields for users to apply any discounts, delivery charges, and additional fees to the total.
* **Calculation and Rounding**: Display the calculated amount each participant owes, with an option to round up each participant's amount to the nearest thousand (or a specified rounding threshold).
* **Summary View**: Show a summary of each participant's items, their individual total, the applied discount, and the final rounded total.
* **Notification/Sharing**: Provide a shareable link or email feature so users can notify participants of their final amounts.

### Backend

* **API for Calculations**: Develop an API endpoint that receives item data, participant assignments, discounts, and fees. It should calculate the total for each participant, apply the discount, and round up the results.
* **OCR Integration**: If implementing OCR for bill uploads, use a service like Google Cloud Vision or Tesseract OCR to extract text from uploaded images.
* **Database**: Use a database (e.g., PostgreSQL, MongoDB) to store:
+ Users' bill information
+ Participant assignments
+ Calculations and historical bills for future reference
* **Authentication**: Allow users to create accounts, save bills, and view their history. Authentication can be done with JWT or OAuth.

### Calculation Logic

* Parse the uploaded bill and assign each item's cost based on user input
* Calculate the total cost of each participant's items
* Apply any discounts and fees proportionally to each participant
* Round each participant's final amount to the nearest thousand or specified rounding value
* Generate a summary of each participant's final share, including rounded values

## Technology Stack

* **Frontend**: React, Vue, or Angular for dynamic UI. Use Tailwind CSS or Bootstrap for styling.
* **Backend**: Node.js with Express or Python with Flask/Django.
* **OCR Service**: Google Cloud Vision API or Tesseract.js for extracting text from images.
* **Database**: PostgreSQL or MongoDB for data storage.
* **Authentication**: JWT for token-based authentication or OAuth for third-party login.
* **Deployment**: Docker for containerization, deployed on AWS, Heroku, or DigitalOcean.

## Example User Workflow

1. User logs in and uploads an image of their bill or inputs item details manually.
2. OCR extracts text (if image uploaded), and user can confirm or edit extracted items.
3. User assigns items to participants, applies any discounts, and submits.
4. Backend processes the data, calculates each participant's share, applies discounts proportionally, and rounds up.
5. Results are displayed with each participant's items, total, discount, and rounded amount.
6. User shares the result via a link or email with participants.

## API Structure

* **POST /api/upload-bill**: Accepts an image or text input, parses items, and returns item names and prices.
* **POST /api/calculate**: Accepts items, participants, discount, and fees to compute the total cost for each participant.
* **POST /api/round-up**: Receives calculated amounts and applies rounding rules.
* **POST /api/share**: Sends the final amounts to each participant by email or generates a shareable link.

## Example User Stories

* As a user, I want to upload my restaurant bill and have the system extract items and prices, so I don't need to input them manually.
* As a user, I want to assign items to participants and see each person's total cost with discounts applied.
* As a user, I want the system to automatically round up each participant's final amount for simplicity.
* As a user, I want to share the final breakdown with my friends through a link or email.

## Implementation Notes

* **OCR**: Use OCR for bill recognition but allow manual correction if the OCR isn't perfect.
* **Rounding Logic**: Customize rounding behavior in the backend to allow users to set their rounding preferences.
* **Discount Handling**: Distribute discounts proportionally based on individual totals to ensure fairness.
