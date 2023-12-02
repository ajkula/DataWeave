-- Enable 'sa' login
ALTER LOGIN sa ENABLE;
-- Reset 'sa' password
ALTER LOGIN sa WITH PASSWORD = 'YourNewStrongPassword!';