import { prefixes } from "./util/id";

export namespace Examples {
  export const Id = (prefix: keyof typeof prefixes) =>
    `${prefixes[prefix]}_XXXXXXXXXXXXXXXXXXXXXXXXX`;

  export const Project = {
    id: Id("project"),
    name: "Spring Boot & React Full-Stack Deployment with Docker",
    year: "2025",
    technologies: "Spring Boot, React, Docker, Docker Compose, Nginx, GitHub Actions, Linux",
    description: "• Containerized both Spring Boot backend and React frontend using multi-stage Docker builds.\n• Deployed the full-stack system using Docker Compose with isolated containers.",
    order: 1,
  };

  export const User = {
    id: Id("user"),
    name: "John Doe",
    email: "john@example.com",
    fingerprint: "183ded44-24d0-480e-9908-c022eff8d111",
    stripeCustomerID: "cus_XXXXXXXXXXXXXXXXX",
  };

  export const Profile = {
    user: User,
  };

  export const Token = {
    id: Id("apiPersonal"),
    token: "trm_test_******XXXX",
    created: new Date("2024-06-29 19:36:19.000"),
  };

  export const App = {
    id: Id("apiClient"),
    secret: "sec_******XXXX",
    name: "Example App",
    redirectURI: "https://example.com/callback",
  };
}
