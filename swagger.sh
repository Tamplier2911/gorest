cd internal \
&& rm -r docs \
&& cd ../cmd/app \
&& swag init --parseDependency true \
&& mv docs ../../internal
