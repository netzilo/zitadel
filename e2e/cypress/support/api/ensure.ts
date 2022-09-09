import { apiCallProperties } from "./apiauth"

// Entity is an object but not a function
type Entity =
    { [k: string]: any }
    & ({ bind?: never }
        | { call?: never });

export function ensureSomething(api: apiCallProperties, searchPath: string, find: (entity: Entity) => boolean, apiPath: (entity: Entity) => string, method: string, body: Entity, expectEntity: (entity: Entity) => boolean): Cypress.Chainable<number> {
    return searchSomething(api, searchPath, find).then(sRes => {
        if (expectEntity(sRes.entity)) {
            return cy.wrap({
                id: sRes?.entity?.id,
                initialSequence: 0
            })
        }

        const req = {
            method: method,
            url: `${api.mgntBaseURL}${apiPath(sRes.entity)}`,
            headers: {
                Authorization: api.authHeader
            },
            body: body,
            failOnStatusCode: false,
            followRedirect: false,
        }

        return cy.request(req).then(cRes => {
            expect(cRes.status).to.equal(200)
            return {
                id: cRes.body.id,
                initialSequence: sRes.sequence
            }
        })
    }).then((data) => {
        awaitDesired(30, expectEntity, data.initialSequence, api, searchPath, find)
        return cy.wrap<number>(data.id)
    })
}

export function ensureSomethingExists(api: apiCallProperties, searchPath: string, find: (entity: Entity) => boolean, createPath: string, body: Entity): Cypress.Chainable<number> {
    return ensureSomething(api, searchPath, find, () => createPath, 'POST', body, entity => !!entity)
}

export function ensureSomethingDoesntExist(api: apiCallProperties, searchPath: string, find: (entity: Entity) => boolean, deletePath: (entity: Entity) => string): Cypress.Chainable<null> {
    return ensureSomething(api, searchPath, find, deletePath, 'DELETE', null, entity => !entity)
        .then(() => { return null })
}

type SearchResult = {
    entity: Entity
    sequence: number
}

export function searchSomething(api: apiCallProperties, searchPath: string, find: (entity: Entity) => boolean): Cypress.Chainable<SearchResult> {

    return cy.request({
        method: 'POST',
        url: `${api.mgntBaseURL}${searchPath}`,
        headers: {
            Authorization: api.authHeader
        },
    }).then(res => {
        return {
            entity: res.body.result?.find(find) || null,
            sequence: res.body.details.processedSequence
        }
    })
}

function awaitDesired(trials: number, expectEntity: (entity: Entity) => boolean, initialSequence: number, api: apiCallProperties, searchPath: string, find: (entity: Entity) => boolean) {
    searchSomething(api, searchPath, find).then(resp => {
        if (!expectEntity(resp.entity) || resp.sequence <= initialSequence) {
            expect(trials, `trying ${trials} more times`).to.be.greaterThan(0);
            cy.wait(1000)
            awaitDesired(trials - 1, expectEntity, initialSequence, api, searchPath, find)
        }
    })
}