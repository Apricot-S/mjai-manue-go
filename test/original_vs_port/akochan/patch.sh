#!/usr/bin/env bash

# Replace removed boost::asio::io_service with boost::asio::io_context
# https://www.boost.org/doc/libs/release/doc/html/boost_asio/history.html
sed -i 's/boost::asio::io_service/boost::asio::io_context/' mjai_client.hpp

# Replace removed boost::asio::buffer_cast with static_cast and data()
# https://www.boost.org/doc/libs/release/doc/html/boost_asio/history.html
sed -i 's/boost::asio::buffer_cast<const char \*>(buffer\.data())/static_cast<const char *>(buffer.data().data())/' mjai_client.cpp

# Change the connection destination from localhost to host.docker.internal
sed -i '
/TcpClient::TcpClient(const std::string &ip, const int port){/,/^}/{
    /socket_->connect(boost::asio::ip::tcp::endpoint(boost::asio::ip::address::from_string(ip), port), error_);/c\
    boost::asio::ip::tcp::resolver resolver(io_service_);\
    boost::asio::ip::tcp::resolver::results_type endpoints = resolver.resolve("host.docker.internal", std::to_string(port));\
    boost::asio::connect(*socket_, endpoints, error_);
}' mjai_client.cpp
