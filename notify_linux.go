package sdnotify

import (
	"fmt"
	"net"
	"os"
)

/*
#include <time.h>
static unsigned long long get_nsecs(void)
{
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (unsigned long long)ts.tv_sec * 1000000000UL + ts.tv_nsec;
}
*/
import "C"

// SdNotify sends a specified string to the systemd notification socket.
func SdNotify(state string) error {
	name := os.Getenv("NOTIFY_SOCKET")
	if name == "" {
		return ErrSdNotifyNoSocket
	}

	conn, err := net.DialUnix("unixgram", nil, &net.UnixAddr{Name: name, Net: "unixgram"})
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(state))
	return err
}

// Reloading sends RELOADING=1\nMONOTONIC_USEC=<monotonic_time> to the 
// systemd notify socket.
func Reloading() error {
	return SdNotify(fmt.Sprintf("RELOADING=1\nMONOTONIC_USEC=%d",
			C.get_nsecs() / 1000))
}
