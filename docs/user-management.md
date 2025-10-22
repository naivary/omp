## Users

This documents describes how it expecets the users of the OMP to be managed and
how they are allowed to interact with different entities.

Right now the following classes of users are defined:

1. Players
2. Coaches
3. Funtionary

The first two calsses are pretty straightforward. Functionaries are people who
are doing work for the club but no metrics are collected for.

Both Players and Coaches can sign up to the platform for themselves
(self-service) but will only be allowed on the platform if the club they have
choosen is accepting them.

Funtionaries can only be created by the root user.

When a Coach joins a Team he unlocks TEAM_OWNER_OPERATIONS.

## Implemetation

The Authetnication will be handled using keycloak as a CNCF certified prodcut
which is open source and supporting industry standard protocols (OIDC, OAuth2,
SAML etc.). Players, Coaches, Root user of the Club (Clubs) and Functionaries
will all be stored in the keycloak to be able to login while their detailed
information will be stored seperate in another database which is managed by the
application. The reference will be set by using custom User attributed and
setting the required id to find the profile. Every token will include the id to
quickly identifie the profile.

I think its even better to store all of the users in one realm e.g. `omp` and
assign the apporpiate roles to the users.
