# Use an official Node.js runtime as a parent image
FROM node:14 AS build

# Set the working directory to /app
WORKDIR /app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the app source code to the working directory
COPY . .

# Build the app for production
RUN npm run build

# Install the serve package globally
RUN npm install -g serve

# Expose port 80
EXPOSE 3000

# Start the application with serve
CMD ["serve", "-s", "build"]

# Use an official Nginx image as a parent image
# FROM nginx:1.21.3-alpine

# # Remove the default Nginx configuration file
# RUN rm /etc/nginx/conf.d/default.conf

# # Copy the Nginx configuration file from the current directory to the container
# COPY ./nginx/nginx.conf /etc/nginx/conf.d/

# # Copy the built app from the previous stage to the Nginx document root directory
# COPY --from=build /app/build /usr/share/nginx/html

# # Expose port 80
# EXPOSE 80

# # Start Nginx in the foreground when the container starts
# CMD ["nginx", "-g", "daemon off;"]
