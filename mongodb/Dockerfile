FROM mongo:8.0

# Create log directory
RUN mkdir -p /var/log/leszmonitor

# MongoDB config
COPY mongod.conf /etc/mongod.conf

# Expose the standard MongoDB port
EXPOSE 27017

# Use the config file
CMD ["mongod", "--config", "/etc/mongod.conf"]