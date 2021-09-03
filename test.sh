# cd internal/v2/comments/tests \
echo "Testing 2nd version of API" \
&& cd internal/v2/auth/tests \
&& go test -v \
&& cd ../../posts/tests \
&& go test -v \
&& cd ../../comments/tests \
&& go test -v \
&& echo "Testing 1st version of API" \
&& cd ../../../v1/posts/tests \
&& go test -v \
&& cd ../../comments/tests \
&& go test -v 