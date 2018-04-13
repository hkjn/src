"""
Implementation of Shamir's secret sharing scheme, based on
https://en.m.wikipedia.org/wiki/Shamir%27s_Secret_Sharing#Python_example.
"""
from __future__ import division, print_function
import random
import functools

# For this application we want a known prime number as close as possible to
# our security level; e.g.  desired security level of 512 bits -- too large and
# all the ciphertext is large; too small and security is compromised).
# 13th Mersenne Prime:
_PRIME = 2**521 - 1


def make_random_shares(secret, minimum=0, shares=0, prime=_PRIME):
    """
    Generates Shamir shares with random coefficients.

    Returns:
        The share points as list of (x, y) valued tuples.
    """

    if minimum > shares:
        raise ValueError("secret would be irrecoverable with minimum > shares")
    if shares <= 1:
        raise ValueError("pool shares must be at least 1")
    if not secret:
        raise ValueError("secret need to be specified")

    def _eval_at(poly, x, prime): # pylint: disable=invalid-name
        """Evaluate polynomial (coefficient list) at x, used to generate a
        shamir pool.

        Args:
            poly: List of coefficients for polynomial. E.g if the polynomial
              is 1234 + 166x + 94x^2, the secret is 1234, and poly will be [1234,
              166, 94].
            x: x value to evaluate polynomial for.
            prime: The prime to use in the prime field evaluation.
        Returns:
            The computed value of the polynomial at x.
        """

        accum = 0
        for coeff in reversed(poly):
            accum *= x
            accum += coeff
            accum %= prime
        return accum

    rint = functools.partial(random.SystemRandom().randint, 0)
    poly = [secret] + [rint(prime) for i in range(1, minimum)]
    print('poly: {}'.format(poly))
    points = [(i, _eval_at(poly, i, prime)) for i in range(1, shares + 1)]
    return points


def _divmod(num, den, p): # pylint: disable=invalid-name
    """Compute num / den modulo prime p.

    To explain what this means, the return value will be such that
    the following is true: den * _divmod(num, den, p) % p == num
    """

    def _extended_gcd(a, b):
        """Calculate greatest common denominator of a and b.

        Division in integers modulus p means finding the inverse of the
        denominator modulo p and then multiplying the numerator by this
        inverse (Note: inverse of A is B such that A*B % p == 1) this can
        be computed via extended Euclidean algorithm
        http://en.wikipedia.org/wiki/Modular_multiplicative_inverse#Computation
        """
        # pylint: disable=invalid-name
        x = 0
        last_x = 1
        while b != 0:
            quot = a // b
            a, b = b, a%b
            x, last_x = last_x - quot * x, x
        return last_x

    inv = _extended_gcd(den, p)
    return num * inv


def _lagrange_interpolate(x, x_s, y_s, p): # pylint: disable=invalid-name
    """Find the y-value for the given x, given n (x, y) points;
    k points will define a polynomial of up to kth order
    """

    if len(x_s) != len(set(x_s)):
        raise ValueError("points must be distinct")

    def _pi(vals):
        """Compute product of input vals.
        """
        accum = 1
        for val in vals:
            accum *= val
        return accum

    nums = [] # avoid inexact division
    dens = []
    for i in range(len(x_s)):
        others = list(x_s)
        cur = others.pop(i)
        nums.append(_pi(x - o for o in others))
        dens.append(_pi(cur - o for o in others))
    den = _pi(dens)
    num = sum(
        (_divmod(nums[i] * den * y_s[i] % p, dens[i], p)
         for i in range(len(x_s)))
    )
    return (_divmod(num, den, p) + p) % p


def recover_secret(shares, prime=_PRIME):
    """Recover the secret from Shamir share points
    (x,y points on the polynomial).
    """

    if len(shares) < 2:
        raise ValueError("need at least two shares")

    x_s, y_s = zip(*shares)
    return _lagrange_interpolate(0, x_s, y_s, prime)


def demo():
    """Demo generates Shamir shares for a static secret and then recovers it
      from those shares.
    """
    secret = 737373
    shares = make_random_shares(secret, minimum=8, shares=10)

    print('secret: {}'.format(secret))
    print('shares: {}'.format(shares))
    print('secret recovered from minimum subset of shares: {}'.format(recover_secret(shares[:8])))
    print('secret recovered from a different minimum subset of shares: {}'.format(
        recover_secret(shares[2:])))


if __name__ == '__main__':
    demo()
