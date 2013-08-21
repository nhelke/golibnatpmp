#include "waitAndReadNATPMPResp.h"

#include <sys/select.h>

int waitAndReadNATPMPResp(natpmp_t *natpmp, resp_s *resp) {
	int r;
	natpmpresp_t response;
	do {
		fd_set fds;
		struct timeval timeout;
		FD_ZERO(&fds);
		FD_SET(natpmp->s, &fds);
		getnatpmprequesttimeout(natpmp, &timeout);
		select(FD_SETSIZE, &fds, NULL, NULL, &timeout);
		r = readnatpmpresponseorretry(natpmp, &response);
	} while (r==NATPMP_TRYAGAIN);

	resp->addr = *((in_addr_t*)(&response.pnu.publicaddress.addr));
	resp->privateport = response.pnu.newportmapping.privateport;
	resp->mappedpublicport = response.pnu.newportmapping.mappedpublicport;
	resp->lifetime = response.pnu.newportmapping.lifetime;

	return r;
}
