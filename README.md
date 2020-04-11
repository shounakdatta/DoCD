<p align="center">
  <img src='./docs/media/DoCD-logo.png?raw=true' width='70%'>
</p>

---

DoCD (pronounced 'Docked') is a containerization alternative to Docker. DoCD uses system relevant package managers for dependency management and has built-in continuous-deployment support with git webhooks.

DoCD is ideal for projects with multiple service dependencies in the same repo
(i.e. React + Flask + MongoDB etc.), developed in machines not powerful enough
to run Docker and a dedicated VM.

It is also useful in environments where implementation of continuous-deployment tools (i.e. Jenkins) is not practical and the developer needs a quick
way to automatically deploy code to development/production servers.
