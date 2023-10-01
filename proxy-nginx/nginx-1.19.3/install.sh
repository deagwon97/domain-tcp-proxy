export LUAJIT_LIB=/usr/lib/x86_64-linux-gnu
export LUAJIT_INC=/usr/include/luajit-2.1

# Here we assume Nginx is to be installed under /opt/nginx/.
./configure --prefix=/opt/nginx \
        --with-ld-opt="-Wl,-rpath,$LUAJIT_LIB" \
        --add-module=/workdir/proxy-nginx/ngx_devel_kit-0.3.2 \
        --add-module=/workdir/proxy-nginx/lua-nginx-module-0.10.25

# Note that you may also want to add `./configure` options which are used in your
# current nginx build.
# You can get usually those options using command nginx -V

# you can change the parallelism number 2 below to fit the number of spare CPU cores in your
# machine.
make -j2
make install

# Note that this version of lug-nginx-module not allow to set `lua_load_resty_core off;` any more.
# So, you have to install `lua-resty-core` and `lua-resty-lrucache` manually as below.

cd lua-resty-core
make install PREFIX=/opt/nginx
cd ..
cd lua-resty-lrucache
make install PREFIX=/opt/nginx
cd ..
cd lua-resty-string
make install PREFIX=/opt/nginx
cd ..
# make

mkdir -p /usr/lib/x86_64-linux-gnu/lua/5.1
mkdir -p /usr/local/lib/lua/5.1/
mkdir -p /usr/local/share/lua/5.1/

ln -s /opt/nginx/lib/lua/* /usr/local/share/lua/5.1/
ln -s /opt/nginx/lib/lua/*/* /usr/lib/x86_64-linux-gnu/lua/5.1
ln -s /opt/nginx/lib/lua/*/* /usr/local/lib/lua/5.1/
ln -s /usr/lib/x86_64-linux-gnu/lua/* /usr/lib/x86_64-linux-gnu/lua/5.1
ln -s /workdir/proxy-nginx/nginx-1.19.3/lua-resty-nettle/lib/resty/* /usr/lib/x86_64-linux-gnu/lua/5.1/resty
ln -s /workdir/proxy-nginx/nginx-1.19.3/lua-resty-nettle/lib/resty/* /usr/local/lib/lua/5.1/resty
# ln -s /workdir/proxy-nginx/nginx-1.19.3/lua-blowfish/build/blowfish.so /usr/local/lib/lua/5.1/
# ln -s /workdir/proxy-nginx/nginx-1.19.3/lua-blowfish/build/blowfish.so /usr/lib/x86_64-linux-gnu/lua/5.1

# git config --global url.https://github.com/.insteadOf git://github.com/

rm -rf /opt/nginx/conf && ln -s /workdir/proxy-nginx/nginx-conf/conf /opt/nginx/conf


# mkdir -p /usr/local/lib/lua/5.1/
# mv /usr/local/lib/lua/resty /usr/local/lib/lua/5.1/
# mv /usr/local/lib/lua/ngx /usr/local/lib/lua/5.1/

# rm -rf /opt/nginx/conf && ln -s /workdir/proxy-nginx/nginx-conf/conf /opt/nginx/conf

 # add necessary `lua_package_path` directive to `nginx.conf`, in the http context

#  lua_package_path "/opt/nginx/lib/lua/?.lua;;";