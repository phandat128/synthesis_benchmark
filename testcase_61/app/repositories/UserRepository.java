package repositories;

import models.UserProfile;
import io.ebean.DB;
import javax.inject.Inject;
import javax.inject.Singleton;
import java.util.Optional;

/**
 * Repository layer for interacting with the UserProfile data store.
 * Ensures secure data access using ORM methods.
 */
@Singleton
public class UserRepository {

    /**
     * Securely retrieves a UserProfile by its primary key ID.
     * This method uses Ebean's built-in find methods, which automatically use parameterized
     * queries internally, making it inherently safe against SQL Injection (CWE-89).
     *
     * @param id The validated numeric ID of the user.
     * @return An Optional containing the UserProfile if found.
     */
    public Optional<UserProfile> findById(Long id) {
        // Using DB.find(Class, ID) is the safest and most efficient way to query by primary key.
        UserProfile profile = DB.find(UserProfile.class, id);
        return Optional.ofNullable(profile);
    }
}