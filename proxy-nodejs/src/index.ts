import net from "net";
import { parseAppAuthority } from "./lib";

require("events").EventEmitter.defaultMaxListeners = 0;

function createTimeout(timeout: number, onTimeout: () => void) {
  let timeoutId: NodeJS.Timeout;
  const refresh = () => {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(onTimeout, timeout);
  };
  const clear = () => clearTimeout(timeoutId);
  return { refresh, clear };
}
const appSockets: { [_: number]: net.Socket } = {};
const userSockets: { [_: number]: net.Socket } = {};

// midServer listen to application
function createAppSocket(
  id: number,
  userSocket: net.Socket,
  applicationHost: string,
  applicationPortNumber: number
) {
  const appSocket = net.createConnection(
    applicationPortNumber,
    applicationHost
  );
  appSockets[id] = appSocket;

  // Clear socket with no data for 5 minutes
  const { refresh, clear } = createTimeout(5 * 60 * 1000, () => {
    appSocket.destroy();
    delete userSockets[id];
  });

  appSocket.on("data", (data) => {
    userSocket.write(data);
    refresh();
  });

  appSocket.on("close", () => {
    userSocket.destroy();
    clear();
    delete userSockets[id];
  });

  appSocket.on("error", (e) => {
    userSocket.destroy();
    clear();
    delete userSockets[id];
  });
}

let counter = 0;
function createUid() {
  counter = (counter + 1) % 10000;
  return counter;
}

function showNextUid() {
  return (counter + 1) % 10000;
}

const midServer = net.createServer(async (midSocket) => {
  // create new tcp connection with new user
  let uid = createUid();

  midSocket.on("data", (data) => {
    // if midSocket receive data from user
    // console.log(`new connect ${uid}`);
    // if first connection
    if (appSockets[uid] === undefined) {
      // if the next uid is already used, destroy it
      const nextUid = showNextUid();
      if (appSockets[nextUid] !== undefined) {
        appSockets[nextUid].destroy();
        delete appSockets[nextUid];
        userSockets[nextUid].destroy();
        delete userSockets[nextUid];
      }

      const { appHost, appPort, err } = parseAppAuthority(data);
      if (err) {
        midSocket.destroy();
        return;
      }
      // connect app to user
      createAppSocket(uid, midSocket, appHost, appPort);
    }

    // if not first connection
    userSockets[uid] = midSocket;

    // if midSocket receive disconnect message from user
    const onDisconnect = (uid: number) => {
      if (appSockets[uid]) {
        appSockets[uid].destroy();
      }
      delete appSockets[uid];
      delete userSockets[uid];
    };
    midSocket.on("close", () => onDisconnect(uid));
    midSocket.on("error", () => onDisconnect(uid));
    midSocket.on("disconnect", () => onDisconnect(uid));

    // just send data to app
    // console.log(`forward data ${uid}`);
    appSockets[uid].write(data);
  });
});

midServer.listen(9980, () => {
  console.log(`run mid server 9980`);
});
